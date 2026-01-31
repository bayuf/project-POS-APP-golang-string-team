package repository

import (
	"context"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	FindAllReservations(ctx context.Context, page, limit int) ([]entity.Reservation, int64, error)
	Create(ctx context.Context, res *entity.Reservation) error
	FindByID(ctx context.Context, id string) (entity.Reservation, error)
	Update(ctx context.Context, res *entity.Reservation) error

	IsTableAvailable(ctx context.Context, tableID int64, resTime time.Time) (bool, error)
}

type reservationRepo struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewReservationRepository(db *gorm.DB, log *zap.Logger) ReservationRepository {
	return &reservationRepo{
		DB:  db,
		Log: log,
	}
}

func (r *reservationRepo) FindAllReservations(ctx context.Context, page, limit int) ([]entity.Reservation, int64, error) {
	var reservations []entity.Reservation
	var total int64
	offset := (page - 1) * limit

	if err := r.DB.WithContext(ctx).Model(&entity.Reservation{}).Count(&total).Error; err != nil {
		r.Log.Error("failed to count reservations", zap.Error(err))
		return nil, 0, err
	}

	if total == 0 {
		return []entity.Reservation{}, 0, nil
	}

	err := r.DB.WithContext(ctx).Debug().
		Model(&entity.Reservation{}).
		Preload("Table").
		Limit(limit).
		Offset(offset).
		Order("reservation_time ASC").
		Find(&reservations).Error

	if err != nil {
		r.Log.Error("failed to find reservations", zap.Error(err))
		return nil, 0, err
	}

	return reservations, total, nil
}

func (r *reservationRepo) Create(ctx context.Context, res *entity.Reservation) error {
	return r.DB.WithContext(ctx).Create(res).Error
}

func (r *reservationRepo) FindByID(ctx context.Context, id string) (entity.Reservation, error) {
	var res entity.Reservation
	err := r.DB.WithContext(ctx).Preload("Table").Where("id = ?", id).First(&res).Error
	return res, err
}

func (r *reservationRepo) Update(ctx context.Context, res *entity.Reservation) error {
	return r.DB.WithContext(ctx).Model(res).
		Select("table_id", "is_cancelled", "updated_at").
		Updates(res).Error
}

func (r *reservationRepo) IsTableAvailable(ctx context.Context, tableID int64, resTime time.Time) (bool, error) {
	var count int64
	// mengecek apakah ada reservasi di meja yang sama dalam rentang 2 jam
	err := r.DB.WithContext(ctx).Model(&entity.Reservation{}).
		Where("table_id = ?", tableID).
		Where("is_cancelled = ?", false).
		Where("reservation_time BETWEEN ? AND ?",
			resTime.Add(-2*time.Hour),
			resTime.Add(2*time.Hour)).
		Count(&count).Error

	return count == 0, err
}
