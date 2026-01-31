package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthService struct {
	repo     repository.AuthRepositoryIface
	logger   *zap.Logger
	tx       *gorm.DB
	emailJob chan<- utils.EmailJob
}

func NewAuthService(repo repository.AuthRepositoryIface, logger *zap.Logger, tx *gorm.DB, emailJob chan<- utils.EmailJob) *AuthService {
	return &AuthService{
		repo:     repo,
		logger:   logger,
		tx:       tx,
		emailJob: emailJob,
	}
}

func (uc *AuthService) Login(ctx context.Context, data dto.Login) (*dto.Session, error) {
	// find user by email
	user, err := uc.repo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		uc.logger.Error("error getting user", zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	// password check
	if !utils.CheckString(user.PasswordHash, data.Password) {
		uc.logger.Error("invalid password")
		return nil, errors.New("invalid credentials")
	}

	// Revoke old session if still active
	if err := uc.repo.RevokeSessionByUserId(ctx, user.ID); err != nil {
		return nil, err
	}

	// Create Session
	id, err := uc.repo.CreateSession(ctx, entity.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	// get latest session
	session, err := uc.repo.GetSession(ctx, *id)
	if err != nil {
		return nil, err
	}

	return &dto.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiredAt,
	}, nil
}

func (uc *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return uc.repo.RevokeSessionBySessionId(ctx, sessionID)
}

func (uc *AuthService) ResetPassword(ctx context.Context, email string) error {
	// get user
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// cek user valid
	if user == nil {
		return errors.New("user not found")
	}

	// generate code otp
	code, err := utils.GenerateOTP()
	if err != nil {
		return err
	}

	hashedCode, err := utils.HashString(code)
	if err != nil {
		return err
	}

	// add OTP to db
	idOTP := uuid.New()
	if err := uc.repo.AddOTP(ctx, entity.OTPRequest{
		ID:        idOTP,
		UserID:    user.ID,
		OTPHash:   hashedCode,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}); err != nil {
		return err
	}

	// send otp via email
	payload := &dto.Email{
		Type:     "otp",
		Email:    user.Email,
		Username: user.Name,
		Code:     code,
	}

	select {
	case uc.emailJob <- utils.EmailJob{Payload: payload}:
	default:
		uc.logger.Warn("email job queue full, skipping email",
			zap.String("email", user.Email),
		)
	}

	return nil
}

func (uc *AuthService) VerifyOTP(ctx context.Context, otp dto.VerifyOTP) (*dto.CodeOTP, error) {
	// get user
	user, err := uc.repo.GetUserByEmail(ctx, otp.Email)
	if err != nil {
		return nil, err
	}

	// get OTP Data
	otpData, err := uc.repo.GetOTPByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// cek OTP
	if !utils.CheckString(otpData.OTPHash, otp.OTPCode) {
		return nil, errors.New("invalid otp")
	}

	return &dto.CodeOTP{
		OTPToken:  &otpData.ID,
		ExpiredAt: otpData.ExpiredAt,
	}, nil
}

func (uc *AuthService) UpdateUserPassword(ctx context.Context, pass dto.UpdatePassword) error {
	// get user
	userOTP, err := uc.repo.GetOTPByID(ctx, pass.Token)
	if err != nil {
		return err
	}

	// cek user valid
	if userOTP == nil {
		return errors.New("token OTP invalid")
	}

	// generate password hash
	hashedPassword, err := utils.HashString(pass.ConfirmPassword)
	if err != nil {
		return err
	}

	// Begin Transactions
	if err = uc.tx.Transaction(func(tx *gorm.DB) error {
		// update user password
		if err := uc.repo.UpdatePasswordUser(ctx, tx, entity.User{
			ID:           userOTP.UserID,
			PasswordHash: hashedPassword,
		}); err != nil {
			return err
		}

		// update OTP status
		if err := uc.repo.UpdateOTPStatus(ctx, tx, userOTP.UserID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
