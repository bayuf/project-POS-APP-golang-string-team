package repository

import (
	"context"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepositoryIface interface {
	GetSession(ctx context.Context, sessionID uuid.UUID) (*entity.Session, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateSession(ctx context.Context, newSession entity.Session) (*uuid.UUID, error)
	RevokeSessionByUserId(ctx context.Context, userID uuid.UUID) error
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

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := entity.User{}
	if err := r.db.WithContext(ctx).
		Model(&user).
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
		Model(&session).
		Where("id = ?", sessionID).
		Where("revoked_at IS NULL").
		Where("used_at IS NULL").
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
		Where("used_at IS NULL").
		Where("revoked_at IS NULL").
		Where("expired_at > ?", time.Now()).
		Update("revoked_at", time.Now()).Error; err != nil {
		r.logger.Error("failed to revoke session", zap.Error(err))
		return err
	}

	return nil
}
