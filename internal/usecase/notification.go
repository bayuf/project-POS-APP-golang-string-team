package usecase

import (
	"context"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationService struct {
	repo   repository.NotificationRepository
	logger *zap.Logger
}

func NewNotificationService(repo repository.NotificationRepository, logg *zap.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		logger: logg,
	}
}

func (uc *NotificationService) GetMyNotifications(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {
	return uc.repo.FindByUserID(ctx, userID)
}

func (uc *NotificationService) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return uc.repo.UpdateStatus(ctx, id, userID, "read")
}

func (uc *NotificationService) RemoveNotification(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return uc.repo.Delete(ctx, id, userID)
}
