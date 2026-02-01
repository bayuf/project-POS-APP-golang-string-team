package repository

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Create(ctx context.Context, n *entity.Notification) error
}

type notifRepo struct {
	DB   *gorm.DB
	Logg *zap.Logger
}

func NewNotificationRepository(db *gorm.DB, logg *zap.Logger) NotificationRepository {
	return &notifRepo{
		DB:   db,
		Logg: logg,
	}
}

func (r *notifRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {
	var notifications []entity.Notification
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		// Preload("User", func(db *gorm.DB) *gorm.DB {
		// 	return db.Select("id", "name", "email", "role")
		// }).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notifRepo) UpdateStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status string) error {
	return r.DB.WithContext(ctx).Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status).Error
}

func (r *notifRepo) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity.Notification{}).Error
}

func (r *notifRepo) Create(ctx context.Context, n *entity.Notification) error {
	return r.DB.WithContext(ctx).Create(n).Error
}
