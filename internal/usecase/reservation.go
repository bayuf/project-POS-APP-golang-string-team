package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationService struct {
	repo   repository.ReservationRepository
	logger *zap.Logger
	tx     *gorm.DB
}

func NewReservationService(repo repository.ReservationRepository, logger *zap.Logger, tx *gorm.DB) *ReservationService {
	return &ReservationService{
		repo:   repo,
		logger: logger,
		tx:     tx,
	}
}

func (u *ReservationService) FindAllReservation(ctx context.Context, page, limit int) ([]entity.Reservation, dto.Pagination, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	reservations, total, err := u.repo.FindAllReservations(ctx, page, limit)
	if err != nil {
		u.logger.Error("failed to get reservations", zap.Error(err))
		return nil, dto.Pagination{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}

	return reservations, pagination, nil
}

func (u *ReservationService) CreateReservation(ctx context.Context, req dto.CreateReservationRequest) (entity.Reservation, error) {
	resTime, err := time.Parse(time.RFC3339, req.ReservationTime)
	if err != nil {
		return entity.Reservation{}, fmt.Errorf("wrong time format (use ISO8601)")
	}

	// reservasi minimal 1 jam dari sekarang
	if resTime.Before(time.Now().Add(1 * time.Hour)) {
		return entity.Reservation{}, fmt.Errorf("Reservations must be made at least 1 hour in advance")
	}

	// ketersediaan meja
	available, err := u.repo.IsTableAvailable(ctx, req.TableID, resTime)
	if err != nil {
		return entity.Reservation{}, err
	}
	if !available {
		return entity.Reservation{}, fmt.Errorf("the table has been booked at that time")
	}

	newID := uuid.New()
	res := entity.Reservation{
		ID:              newID,
		CustomerName:    req.CustomerName,
		TableID:         req.TableID,
		ReservationTime: resTime,
		IsCancelled:     false,
	}

	u.logger.Info("load", zap.String("customer", req.CustomerName))

	if err := u.repo.Create(ctx, &res); err != nil {
		u.logger.Error("Error create reservation", zap.Error(err))
		return entity.Reservation{}, err
	}

	return u.repo.FindByID(ctx, newID.String())
}

func (u *ReservationService) UpdateReservation(ctx context.Context, id string, req dto.UpdateReservationRequest) (entity.Reservation, error) {
	existing, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return entity.Reservation{}, err
	}

	now := time.Now()
	existing.TableID = req.TableID
	existing.IsCancelled = req.IsCancelled
	existing.UpdatedAt = &now

	if err := u.repo.Update(ctx, &existing); err != nil {
		return entity.Reservation{}, err
	}

	return existing, nil
}

func (u *ReservationService) GetReservationDetail(ctx context.Context, id string) (entity.Reservation, error) {
	return u.repo.FindByID(ctx, id)
}
