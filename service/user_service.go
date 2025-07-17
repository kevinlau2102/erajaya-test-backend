package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"erajaya-interview/constants"
	"erajaya-interview/dto"
	"erajaya-interview/entity"
	"erajaya-interview/helpers"
	"erajaya-interview/repository"
)

type (
	UserService interface {
		Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req dto.PaginationRequest) (dto.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
		GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error)
		Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (dto.UserUpdateResponse, error)
		Delete(ctx context.Context, userId string) error
		Verify(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error)
		RevokeRefreshToken(ctx context.Context, userID string) error
	}

	userService struct {
		userRepo         repository.UserRepository
		refreshTokenRepo repository.RefreshTokenRepository
		jwtService       JWTService
		db               *gorm.DB
	}
)

func NewUserService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService JWTService,
	db *gorm.DB,
) UserService {
	return &userService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		db:               db,
	}
}

func SafeRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		// TODO: Do you think that we should panic here?
		// panic(r)
	}
}

func (s *userService) Register(ctx context.Context, req dto.UserCreateRequest) (dto.UserResponse, error) {

	_, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserResponse{}, err
	}

	if flag {
		return dto.UserResponse{}, dto.ErrEmailAlreadyExists
	}

	user := entity.User{
		Name:       req.Name,
		TelpNumber: req.TelpNumber,
		Role:       constants.ENUM_ROLE_USER,
		Email:      req.Email,
		Password:   req.Password,
	}

	userReg, err := s.userRepo.Register(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, dto.ErrCreateUser
	}

	return dto.UserResponse{
		ID:         userReg.ID.String(),
		Name:       userReg.Name,
		TelpNumber: userReg.TelpNumber,
		Role:       userReg.Role,
		Email:      userReg.Email,
	}, nil
}

func (s *userService) GetAllUserWithPagination(
	ctx context.Context,
	req dto.PaginationRequest,
) (dto.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return dto.UserPaginationResponse{}, err
	}

	var datas []dto.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := dto.UserResponse{
			ID:         user.ID.String(),
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
			TelpNumber: user.TelpNumber,
		}

		datas = append(datas, data)
	}

	return dto.UserPaginationResponse{
		Data: datas,
		PaginationResponse: dto.PaginationResponse{
			Page:    dataWithPaginate.Page,
			Limit:   dataWithPaginate.Limit,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserById
	}

	return dto.UserResponse{
		ID:         user.ID.String(),
		Name:       user.Name,
		TelpNumber: user.TelpNumber,
		Role:       user.Role,
		Email:      user.Email,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (dto.UserResponse, error) {
	emails, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return dto.UserResponse{}, dto.ErrGetUserByEmail
	}

	return dto.UserResponse{
		ID:         emails.ID.String(),
		Name:       emails.Name,
		TelpNumber: emails.TelpNumber,
		Role:       emails.Role,
		Email:      emails.Email,
	}, nil
}

func (s *userService) Update(ctx context.Context, req dto.UserUpdateRequest, userId string) (
	dto.UserUpdateResponse,
	error,
) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	data := entity.User{
		ID:         user.ID,
		Name:       req.Name,
		TelpNumber: req.TelpNumber,
		Role:       user.Role,
		Email:      req.Email,
	}

	userUpdate, err := s.userRepo.Update(ctx, nil, data)
	if err != nil {
		return dto.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return dto.UserUpdateResponse{
		ID:         userUpdate.ID.String(),
		Name:       userUpdate.Name,
		TelpNumber: userUpdate.TelpNumber,
		Role:       userUpdate.Role,
		Email:      userUpdate.Email,
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	err = s.userRepo.Delete(ctx, nil, user.ID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (s *userService) Verify(ctx context.Context, req dto.UserLoginRequest) (dto.TokenResponse, error) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	user, err := s.userRepo.GetUserByEmail(ctx, tx, req.Email)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	checkPassword, err := helpers.CheckPassword(user.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New("invalid email or password")
	}

	accessToken := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)

	refreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	hashedToken, err := helpers.HashPassword(refreshTokenString)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, user.ID.String()); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	refreshToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}

	if _, err := s.refreshTokenRepo.Create(ctx, tx, refreshToken); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		Role:         user.Role,
	}, nil
}

func (s *userService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.TokenResponse, error) {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	// Find the refresh token in the database
	dbToken, err := s.refreshTokenRepo.FindByPlainToken(ctx, tx, req.RefreshToken)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New(dto.MESSAGE_FAILED_INVALID_REFRESH_TOKEN)
	}

	if time.Now().After(dbToken.ExpiresAt) {
		tx.Rollback()
		return dto.TokenResponse{}, errors.New(dto.MESSAGE_FAILED_EXPIRED_REFRESH_TOKEN)
	}

	user, err := s.userRepo.GetUserById(ctx, tx, dbToken.UserID.String())
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, dto.ErrUserNotFound
	}

	accessToken := s.jwtService.GenerateAccessToken(user.ID.String(), user.Role)

	refreshTokenString, expiresAt := s.jwtService.GenerateRefreshToken()

	hashedToken, err := helpers.HashPassword(refreshTokenString)
	if err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, user.ID.String()); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	refreshToken := entity.RefreshToken{
		UserID:    user.ID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}

	if _, err := s.refreshTokenRepo.Create(ctx, tx, refreshToken); err != nil {
		tx.Rollback()
		return dto.TokenResponse{}, err
	}

	if err := tx.Commit().Error; err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		Role:         user.Role,
	}, nil
}

func (s *userService) RevokeRefreshToken(ctx context.Context, userID string) error {
	tx := s.db.Begin()
	defer SafeRollback(tx)

	// Check if user exists
	_, err := s.userRepo.GetUserById(ctx, tx, userID)
	if err != nil {
		tx.Rollback()
		return dto.ErrUserNotFound
	}

	// Delete all refresh tokens for the user
	if err := s.refreshTokenRepo.DeleteByUserID(ctx, tx, userID); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
