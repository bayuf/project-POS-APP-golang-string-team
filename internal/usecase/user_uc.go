package usecase

import (
	"context"

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

	shiftStart, err := utils.ParseTime(newUser.ShiftStart)
	if err != nil {
		return err
	}

	shiftEnd, err := utils.ParseTime(newUser.ShiftEnd)
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
		ShiftStart:       shiftStart,
		ShiftEnd:         shiftEnd,
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

	shiftStart, err := utils.ParseTime(newUserData.ShiftStart)
	if err != nil {
		return err
	}

	shiftEnd, err := utils.ParseTime(newUserData.ShiftEnd)
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
		ShiftStart:       shiftStart,
		ShiftEnd:         shiftEnd,
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
