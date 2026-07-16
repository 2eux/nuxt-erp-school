package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/opencode/erp-school-backend/internal/config"
	"github.com/opencode/erp-school-backend/internal/domain"
	"github.com/opencode/erp-school-backend/internal/dto"
	"github.com/opencode/erp-school-backend/internal/infrastructure/cache"
	"github.com/opencode/erp-school-backend/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest, userAgent, ipAddress string) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest, userAgent, ipAddress string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error
	GetMe(ctx context.Context, userID string) (*dto.UserDetail, error)
	UpdateProfile(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserBrief, error)
	ValidateAccessToken(tokenString string) (*Claims, error)
	ValidateRefreshToken(tokenString string) (string, error)
	GenerateAccessToken(userID, schoolID, email string, roles []string) (string, time.Time, error)
	GenerateRefreshToken(userID string) (string, time.Time, error)
}

type Claims struct {
	UserID   string   `json:"user_id"`
	SchoolID string   `json:"school_id"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

type authService struct {
	authRepo    repository.AuthRepository
	schoolRepo  repository.SchoolRepository
	redisCache  *cache.RedisClient
	cfg         config.JWTConfig
	logger      *zap.Logger
}

func NewAuthService(
	authRepo repository.AuthRepository,
	schoolRepo repository.SchoolRepository,
	redisCache *cache.RedisClient,
	cfg config.JWTConfig,
	logger *zap.Logger,
) AuthService {
	return &authService{
		authRepo:   authRepo,
		schoolRepo: schoolRepo,
		redisCache: redisCache,
		cfg:        cfg,
		logger:     logger,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, userAgent, ipAddress string) (*dto.LoginResponse, error) {
	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if !user.IsActive {
		return nil, domain.NewUnauthorizedError("user is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	roles, err := s.authRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to get user roles", err)
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Slug
	}

	accessToken, accessExp, err := s.GenerateAccessToken(user.ID, user.SchoolID, user.Email, roleNames)
	if err != nil {
		return nil, domain.NewInternalError("failed to generate access token", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to generate refresh token", err)
	}

	session := &domain.UserSession{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    refreshExp,
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, domain.NewInternalError("failed to create session", err)
	}

	if err := s.authRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		s.logger.Warn("failed to update last login", zap.String("user_id", user.ID), zap.Error(err))
	}

	userBrief := dto.UserBrief{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
		Roles:     roleNames,
		SchoolID:  user.SchoolID,
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		User:         userBrief,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest, userAgent, ipAddress string) (*dto.TokenResponse, error) {
	userID, err := s.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, domain.NewUnauthorizedError("invalid refresh token")
	}

	session, err := s.authRepo.FindSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, domain.NewUnauthorizedError("session not found")
	}

	if err := s.authRepo.DeleteSession(ctx, session.ID); err != nil {
		s.logger.Warn("failed to delete old session", zap.Error(err))
	}

	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, domain.NewUnauthorizedError("user not found")
	}

	if !user.IsActive {
		return nil, domain.NewUnauthorizedError("user is inactive")
	}

	roles, err := s.authRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to get roles", err)
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Slug
	}

	accessToken, accessExp, err := s.GenerateAccessToken(user.ID, user.SchoolID, user.Email, roleNames)
	if err != nil {
		return nil, domain.NewInternalError("failed to generate access token", err)
	}

	refreshToken, refreshExp, err := s.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to generate refresh token", err)
	}

	newSession := &domain.UserSession{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		ExpiresAt:    refreshExp,
	}

	if err := s.authRepo.CreateSession(ctx, newSession); err != nil {
		return nil, domain.NewInternalError("failed to create session", err)
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.authRepo.FindSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil
	}
	return s.authRepo.DeleteSession(ctx, session.ID)
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	_, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	resetToken := uuid.New().String()
	resetKey := fmt.Sprintf("reset_password:%s", resetToken)
	if err := s.redisCache.SetString(ctx, resetKey, req.Email, 15*time.Minute); err != nil {
		return domain.NewInternalError("failed to store reset token", err)
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	resetKey := fmt.Sprintf("reset_password:%s", req.Token)
	email, err := s.redisCache.Get(ctx, resetKey)
	if err != nil || email == "" {
		return domain.NewInvalidInputError("invalid or expired reset token")
	}

	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return domain.NewNotFoundError("user", email)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.NewInternalError("failed to hash password", err)
	}

	if err := s.authRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return domain.NewInternalError("failed to update password", err)
	}

	s.redisCache.Delete(ctx, resetKey)

	return s.authRepo.DeleteUserSessions(ctx, user.ID)
}

func (s *authService) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return domain.NewNotFoundError("user", userID)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return domain.NewInvalidInputError("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return domain.NewInternalError("failed to hash password", err)
	}

	if err := s.authRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return domain.NewInternalError("failed to update password", err)
	}

	return s.authRepo.DeleteUserSessions(ctx, user.ID)
}

func (s *authService) GetMe(ctx context.Context, userID string) (*dto.UserDetail, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, domain.NewNotFoundError("user", userID)
	}

	roles, err := s.authRepo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to get roles", err)
	}

	permissions, err := s.authRepo.GetUserPermissions(ctx, user.ID)
	if err != nil {
		return nil, domain.NewInternalError("failed to get permissions", err)
	}

	roleBriefs := make([]dto.RoleBrief, len(roles))
	for i, r := range roles {
		roleBriefs[i] = dto.RoleBrief{
			ID:   r.ID,
			Name: r.Name,
			Slug: r.Slug,
		}
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

func (s *authService) UpdateProfile(ctx context.Context, userID string, req dto.UpdateUserRequest) (*dto.UserBrief, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, domain.NewNotFoundError("user", userID)
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

	if err := s.authRepo.UpdateUser(ctx, user); err != nil {
		return nil, domain.NewInternalError("failed to update profile", err)
	}

	roles, _ := s.authRepo.GetUserRoles(ctx, user.ID)
	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Slug
	}

	return &dto.UserBrief{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		Phone:     user.Phone,
		Roles:     roleNames,
		SchoolID:  user.SchoolID,
	}, nil
}

func (s *authService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	return claims, nil
}

func (s *authService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.RefreshSecret), nil
	})
	if err != nil {
		return "", domain.ErrTokenExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", domain.ErrTokenInvalid
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", domain.ErrTokenInvalid
	}

	return sub, nil
}

func (s *authService) GenerateAccessToken(userID, schoolID, email string, roles []string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.AccessTTL)
	claims := &Claims{
		UserID:   userID,
		SchoolID: schoolID,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.Issuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *authService) GenerateRefreshToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.RefreshTTL)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    s.cfg.Issuer,
		ID:        userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}
