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
	// create uuid
	newUUID := uuid.New()

	// hash Password
	hashedPassword, err := utils.HashString("12345")
	if err != nil {
		uc.logger.Error("failed to hash password", zap.Error(err))
		return err
	}

	if err := uc.repo.CreateUser(ctx, entity.User{
		ID:               newUUID,
		Name:             newUser.Name,
		Email:            newUser.Email,
		Phone:            newUser.Phone,
		BirthDate:        newUser.Birthdate,
		Salary:           newUser.Salary,
		PasswordHash:     hashedPassword,
		Role:             newUser.Role,
		Address:          newUser.Address,
		AdditionalDetail: newUser.AdditionalDetail,
		AvatarURL:        *newUser.AvatarURL,
		ShiftStart:       &newUser.ShiftStart,
		ShiftEnd:         &newUser.ShiftEnd,
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

// func (uc *UserService) UpdateUserData(ctx context.Context, ID uuid.UUID, newUserData dto.) error {
// 	uc.repo.UpdateUserByID(ctx, ID, )
// }
