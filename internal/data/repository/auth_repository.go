package repository

import (
	"context"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepositoryIface interface {
	ValidateSession(ctx context.Context, sessionID uuid.UUID) (*dto.ValidateSession, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateSession(ctx context.Context, newSession entity.Session) (*uuid.UUID, error)
	RevokeSessionByUserId(ctx context.Context, userID uuid.UUID) error
	RevokeSessionBySessionId(ctx context.Context, sessionID uuid.UUID) error
	AddOTP(ctx context.Context, newOTP entity.OTPRequest) error
	GetOTPByID(ctx context.Context, id uuid.UUID) (*entity.OTPRequest, error)
	GetOTPByUserID(ctx context.Context, id uuid.UUID) (*entity.OTPRequest, error)
	UpdateOTPStatus(ctx context.Context, userID uuid.UUID) error
	UpdatePasswordUser(ctx context.Context, user entity.User) error
}

type AuthRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zap.Logger) *AuthRepository {
	return &AuthRepository{
		db:     db,
		logger: logger,
	}
}

func (r *AuthRepository) ValidateSession(ctx context.Context, sessionID uuid.UUID) (*dto.ValidateSession, error) {
	session := entity.Session{}
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("sessions.id = ?", sessionID).
		Where("expired_at > ?", time.Now()).
		Where("revoked_at IS NULL").
		First(&session).Error; err != nil {
		r.logger.Error("failed to validate session", zap.Error(err))
		return nil, err
	}

	return &dto.ValidateSession{
		SessionID: session.ID,
		UserID:    session.UserID,
		Role:      session.User.Role,
	}, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := entity.User{}
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		r.logger.Error("failed to get user by email", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, newSession entity.Session) (*uuid.UUID, error) {

	if err := r.db.WithContext(ctx).
		Model(&newSession).
		Create(&newSession).Error; err != nil {
		r.logger.Error("failed to create session", zap.Error(err))
		return nil, err
	}
	id := newSession.ID

	return &id, nil
}

func (r *AuthRepository) GetSession(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error) {
	session := entity.Session{}
	if err := r.db.WithContext(ctx).
		Where("id = ?", sessionID).
		Where("revoked_at IS NULL").
		Where("expired_at > ?", time.Now()).
		First(&session).Error; err != nil {
		r.logger.Error("failed to get session", zap.Error(err))
		return nil, err
	}

	return &session, nil
}

func (r *AuthRepository) RevokeSessionByUserId(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Session{}).
		Where("user_id = ?", userID).
		Where("revoked_at IS NULL").
		Where("expired_at > ?", time.Now()).
		Update("revoked_at", time.Now()).Error; err != nil {
		r.logger.Error("failed to revoke session", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepository) RevokeSessionBySessionId(ctx context.Context, sessionID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.Session{}).
		Where("id = ?", sessionID).
		Where("revoked_at IS NULL").
		Update("revoked_at", time.Now()).Error; err != nil {
		r.logger.Error("failed to revoke session", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepository) AddOTP(ctx context.Context, newOTP entity.OTPRequest) error {
	if err := r.db.WithContext(ctx).
		Model(&newOTP).
		Create(&newOTP).Error; err != nil {
		r.logger.Error("failed to create otp", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepository) GetOTPByID(ctx context.Context, id uuid.UUID) (*entity.OTPRequest, error) {
	otp := entity.OTPRequest{}
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("expired_at > ?", time.Now()).
		Where("is_used = ?", false).
		First(&otp).Error; err != nil {
		r.logger.Error("failed to get otp", zap.Error(err))
		return nil, err
	}

	return &otp, nil
}

func (r *AuthRepository) GetOTPByUserID(ctx context.Context, id uuid.UUID) (*entity.OTPRequest, error) {
	otp := entity.OTPRequest{}
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", id).
		Where("expired_at > ?", time.Now()).
		Where("is_used = ?", false).
		First(&otp).Error; err != nil {
		r.logger.Error("failed to get otp", zap.Error(err))
		return nil, err
	}

	return &otp, nil
}

func (r *AuthRepository) UpdateOTPStatus(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&entity.OTPRequest{}).
		Where("user_id = ?", userID).
		Where("is_used = ?", false).
		Update("is_used", true).Error; err != nil {
		r.logger.Error("failed to update otp status", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepository) UpdatePasswordUser(ctx context.Context, user entity.User) error {
	if err := r.db.WithContext(ctx).
		Model(&user).
		Where("id = ?", user.ID).
		Where("is_active = ?", true).
		Update("password_hash", user.PasswordHash).
		Error; err != nil {
		r.logger.Error("failed to update password", zap.Error(err))
		return err
	}

	return nil
}
