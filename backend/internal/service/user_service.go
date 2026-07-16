package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	ListUsers(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.UserDetail, int64, error)
	GetUser(ctx context.Context, id string) (*dto.UserDetail, error)
	CreateUser(ctx context.Context, schoolID string, req dto.CreateUserRequest) (*dto.UserDetail, error)
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserDetail, error)
	DeleteUser(ctx context.Context, id string) error
	ListRoles(ctx context.Context, schoolID string) ([]dto.RoleDetail, error)
	CreateRole(ctx context.Context, schoolID string, req dto.CreateRoleRequest) (*dto.RoleDetail, error)
	UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) (*dto.RoleDetail, error)
	ListPermissions(ctx context.Context) ([]dto.PermissionBrief, error)
}

type userService struct {
	authRepo   repository.AuthRepository
	schoolRepo repository.SchoolRepository
	logger     *zap.Logger
}

func NewUserService(authRepo repository.AuthRepository, schoolRepo repository.SchoolRepository, logger *zap.Logger) UserService {
	return &userService{authRepo: authRepo, schoolRepo: schoolRepo, logger: logger}
}

func (s *userService) ListUsers(ctx context.Context, schoolID string, filter dto.PaginationRequest) ([]dto.UserDetail, int64, error) {
	return nil, 0, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserDetail, error) {
	user, err := s.authRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("user", id)
	}

	roles, _ := s.authRepo.GetUserRoles(ctx, user.ID)
	permissions, _ := s.authRepo.GetUserPermissions(ctx, user.ID)

	roleBriefs := make([]dto.RoleBrief, len(roles))
	for i, r := range roles {
		roleBriefs[i] = dto.RoleBrief{ID: r.ID, Name: r.Name, Slug: r.Slug}
	}

	return &dto.UserDetail{
		ID:               user.ID,
		SchoolID:         user.SchoolID,
		Email:            user.Email,
		Username:         user.Username,
		FullName:         user.FullName,
		AvatarURL:        user.AvatarURL,
		Phone:            user.Phone,
		IsActive:         user.IsActive,
		EmailVerifiedAt:  user.EmailVerifiedAt,
		LastLoginAt:      user.LastLoginAt,
		PasswordChangedAt: user.PasswordChangedAt,
		Roles:            roleBriefs,
		Permissions:      permissions,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, schoolID string, req dto.CreateUserRequest) (*dto.UserDetail, error) {
	existing, _ := s.authRepo.FindUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, domain.NewDuplicateError("user", req.Email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewInternalError("failed to hash password", err)
	}

	user := &domain.User{
		ID:           uuid.New().String(),
		SchoolID:     schoolID,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     true,
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, domain.NewInternalError("failed to create user", err)
	}

	if len(req.RoleIDs) > 0 {
		if err := s.authRepo.AssignRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, domain.NewInternalError("failed to assign roles", err)
		}
	}

	return s.GetUser(ctx, user.ID)
}

func (s *userService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserDetail, error) {
	user, err := s.authRepo.FindUserByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("user", id)
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.authRepo.UpdateUser(ctx, user); err != nil {
		return nil, domain.NewInternalError("failed to update user", err)
	}

	if len(req.RoleIDs) > 0 {
		if err := s.authRepo.AssignRoles(ctx, user.ID, req.RoleIDs); err != nil {
			return nil, domain.NewInternalError("failed to update roles", err)
		}
	}

	return s.GetUser(ctx, user.ID)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.authRepo.DeleteUserSessions(ctx, id); err != nil {
		return domain.NewInternalError("failed to delete sessions", err)
	}
	return nil
}

func (s *userService) ListRoles(ctx context.Context, schoolID string) ([]dto.RoleDetail, error) {
	roles, err := s.schoolRepo.ListRoles(ctx, schoolID)
	if err != nil {
		return nil, domain.NewInternalError("failed to list roles", err)
	}

	result := make([]dto.RoleDetail, len(roles))
	for i, r := range roles {
		result[i] = dto.RoleDetail{
			ID:          r.ID,
			SchoolID:    r.SchoolID,
			Name:        r.Name,
			Slug:        r.Slug,
			Description: r.Description,
			IsSystem:    r.IsSystem,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
		}
	}
	return result, nil
}

func (s *userService) CreateRole(ctx context.Context, schoolID string, req dto.CreateRoleRequest) (*dto.RoleDetail, error) {
	slug := uuid.New().String()[:8]
	role := &domain.Role{
		ID:          uuid.New().String(),
		SchoolID:    schoolID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		IsSystem:    false,
	}

	if err := s.schoolRepo.CreateRole(ctx, role, req.Permissions); err != nil {
		return nil, domain.NewInternalError("failed to create role", err)
	}

	return &dto.RoleDetail{
		ID:          role.ID,
		SchoolID:    role.SchoolID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (s *userService) UpdateRole(ctx context.Context, id string, req dto.UpdateRoleRequest) (*dto.RoleDetail, error) {
	role, err := s.schoolRepo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, domain.NewNotFoundError("role", id)
	}

	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	permIDs := req.Permissions
	if len(permIDs) == 0 {
		permIDs = []string{}
	}

	if err := s.schoolRepo.UpdateRole(ctx, role, permIDs); err != nil {
		return nil, domain.NewInternalError("failed to update role", err)
	}

	return &dto.RoleDetail{
		ID:          role.ID,
		SchoolID:    role.SchoolID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}, nil
}

func (s *userService) ListPermissions(ctx context.Context) ([]dto.PermissionBrief, error) {
	perms, err := s.schoolRepo.ListPermissions(ctx)
	if err != nil {
		return nil, domain.NewInternalError("failed to list permissions", err)
	}

	result := make([]dto.PermissionBrief, len(perms))
	for i, p := range perms {
		result[i] = dto.PermissionBrief{ID: p.ID, Name: p.Name, Slug: p.Slug}
	}
	return result, nil
}
