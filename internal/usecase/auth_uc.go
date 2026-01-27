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
)

type AuthService struct {
	repo   repository.AuthRepositoryIface
	logger *zap.Logger
	// tx     *gorm.Tx
}

func NewAuthService(repo repository.AuthRepositoryIface, logger *zap.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
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
