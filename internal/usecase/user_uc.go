package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/entity"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/data/repository"
	"github.com/bayuf/project-POS-APP-golang-string-team/internal/dto"
	"github.com/bayuf/project-POS-APP-golang-string-team/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	repo   repository.UserRepositoryIface
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepositoryIface, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (uc *UserService) CreateUser(ctx context.Context, newUser dto.CreateUser) error {
	hashedPassword, err := utils.HashString("12345")
	if err != nil {
		return err
	}

	birthDate, err := utils.ParseDate(newUser.Birthdate)
	if err != nil {
		return err
	}

	if err := uc.repo.CreateUser(ctx, entity.User{
		ID:               uuid.New(),
		Name:             newUser.Name,
		Email:            newUser.Email,
		Phone:            newUser.Phone,
		BirthDate:        birthDate,
		Salary:           newUser.Salary,
		PasswordHash:     hashedPassword,
		Role:             newUser.Role,
		Address:          newUser.Address,
		AdditionalDetail: newUser.AdditionalDetail,
		AvatarURL:        *newUser.AvatarURL,
		ShiftStart:       newUser.ShiftStart,
		ShiftEnd:         newUser.ShiftEnd,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UserService) GetUserByID(ctx context.Context, ID uuid.UUID) (*dto.UserDetail, error) {
	data, err := uc.repo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserDetail{
		Name:       data.Name,
		Email:      data.Email,
		Phone:      data.Phone,
		BirthDate:  data.BirthDate,
		Role:       data.Role,
		Address:    data.Address,
		Salary:     data.Salary,
		AvatarURL:  data.AvatarURL,
		ShiftStart: data.ShiftStart,
		ShiftEnd:   data.ShiftEnd,
	}, nil
}

func (uc *UserService) UpdateUserData(ctx context.Context, ID uuid.UUID, newUserData dto.UpdateUser) error {
	birthDate, err := utils.ParseDate(newUserData.Birthdate)
	if err != nil {
		return err
	}

	if err := uc.repo.UpdateUserByID(ctx, ID, entity.User{
		Name:             newUserData.Name,
		Email:            newUserData.Email,
		Phone:            newUserData.Phone,
		BirthDate:        birthDate,
		Salary:           newUserData.Salary,
		Role:             newUserData.Role,
		Address:          newUserData.Address,
		AdditionalDetail: newUserData.AdditionalDetail,
		AvatarURL:        *newUserData.AvatarURL,
		ShiftStart:       newUserData.ShiftStart,
		ShiftEnd:         newUserData.ShiftEnd,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UserService) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	if err := uc.repo.DeleteUserByID(ctx, ID); err != nil {
		return err
	}
	return nil
}

func (uc *UserService) GetAllUser(ctx context.Context, req dto.UserFilterRequest) (*[]dto.UserLists, *dto.Pagination, error) {
	// Set default limit jika kosong
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 5
	}

	// 1. Panggil Repository
	users, total, err := uc.repo.ListUser(ctx, req)
	if err != nil {
		return nil, &dto.Pagination{}, err
	}

	var userResponse []dto.UserLists

	for _, t := range *users {
		now := time.Now()
		age := now.Year() - t.BirthDate.Year()

		if now.Month() < t.BirthDate.Month() || (now.Month() == t.BirthDate.Month() && now.Day() < t.BirthDate.Day()) {
			age--
		}

		var shiftStart, shiftEnd string
		var timing string
		if t.ShiftStart != "" || t.ShiftEnd != "" {
			shiftStart = t.ShiftStart
			shiftEnd = t.ShiftEnd

			timing = fmt.Sprintf("%s to %s", shiftStart, shiftEnd)
		}

		res := dto.UserLists{
			ID:        t.ID,
			Name:      t.Name,
			Email:     t.Email,
			Phone:     t.Phone,
			Age:       age,
			Salary:    t.Salary,
			Role:      t.Role,
			AvatarURL: t.AvatarURL,
			Timing:    timing,
		}
		userResponse = append(userResponse, res)
	}

	totalPages := utils.TotalPage(req.Limit, total)

	pagination := dto.Pagination{
		CurrentPage:  req.Page,
		Limit:        req.Limit,
		TotalPages:   totalPages,
		TotalRecords: total,
	}

	return &userResponse, &pagination, nil
}
