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
	repo     repository.UserRepositoryIface
	logger   *zap.Logger
	emailJob chan<- utils.EmailJob
}

func NewUserService(repo repository.UserRepositoryIface, logger *zap.Logger, emailJob chan<- utils.EmailJob) *UserService {
	return &UserService{
		repo:     repo,
		logger:   logger,
		emailJob: emailJob,
	}
}

func (uc *UserService) CreateUser(ctx context.Context, newUser dto.CreateUser) error {
	defaultPass, err := utils.GenerateRandomString(6)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashString(defaultPass)
	if err != nil {
		return err
	}

	birthDate, err := utils.ParseDate(newUser.Birthdate)
	if err != nil {
		return err
	}

	// default permissions
	if newUser.Permissions == nil {
		newUser.Permissions = map[string]bool{
			entity.PermissionDashboard: true,
			entity.PermissionReports:   false,
			entity.PermissionInventory: false,
			entity.PermissionOrders:    false,
			entity.PermissionCustomers: false,
			entity.PermissionSettings:  false,
		}
	}

	permissions := entity.UserPermissions(newUser.Permissions)

	if err := uc.repo.CreateUser(ctx, entity.User{
		ID:               uuid.New(),
		Name:             newUser.Name,
		Email:            newUser.Email,
		Phone:            newUser.Phone,
		BirthDate:        birthDate,
		Salary:           newUser.Salary,
		PasswordHash:     hashedPassword,
		Role:             newUser.Role,
		Permissions:      permissions,
		Address:          newUser.Address,
		AdditionalDetail: newUser.AdditionalDetail,
		AvatarURL:        *newUser.AvatarURL,
		ShiftStart:       newUser.ShiftStart,
		ShiftEnd:         newUser.ShiftEnd,
	}); err != nil {
		return err
	}

	// send password via email
	payload := &dto.Email{
		Type:     "password",
		Email:    newUser.Email,
		Username: newUser.Name,
		Password: defaultPass,
	}

	select {
	case uc.emailJob <- utils.EmailJob{Payload: payload}:
	default:
		uc.logger.Warn("email job queue full, skipping email",
			zap.String("email", newUser.Email),
		)
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

func (uc *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfile, error) {
	user, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		Name:      user.Name,
		Role:      user.Role,
		Email:     user.Email,
		Address:   user.Address,
		AvatarURL: user.AvatarURL,
	}, nil
}

func (uc *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, newUserData dto.UpdateUserProfile) error {
	user, err := uc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	var hashedPassword string
	if newUserData.NewPassword != "" {
		var err error
		hashedPassword, err = utils.HashString(newUserData.NewPassword)
		if err != nil {
			return err
		}
	}

	if err := uc.repo.UpdateUserByID(ctx, userID, entity.User{
		Name:         newUserData.Name,
		Email:        newUserData.Email,
		Address:      newUserData.Address,
		PasswordHash: hashedPassword,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UserService) UpdatePermissions(ctx context.Context, userID uuid.UUID, newPermissions map[string]bool) error {
	user, err := uc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.Role != "admin" {
		return fmt.Errorf("only admin can update permissions")
	}

	if err := uc.repo.UpdateUserByID(ctx, userID, entity.User{
		Permissions: newPermissions,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *UserService) GetAllAdmins(ctx context.Context) ([]dto.ListAdmin, error) {
	admins, err := uc.repo.GetAllAdmins(ctx)
	if err != nil {
		return nil, err
	}

	var listAdmins []dto.ListAdmin
	for _, admin := range *admins {
		listAdmins = append(listAdmins, dto.ListAdmin{
			ID:          admin.ID,
			Name:        admin.Name,
			Email:       admin.Email,
			Permissions: admin.Permissions,
		})
	}

	return listAdmins, nil
}
