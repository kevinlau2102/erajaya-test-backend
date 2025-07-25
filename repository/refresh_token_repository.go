package repository

import (
	"context"
	"time"

	"erajaya-interview/entity"
	"erajaya-interview/helpers"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, tx *gorm.DB, token entity.RefreshToken) (entity.RefreshToken, error)
	FindByToken(ctx context.Context, tx *gorm.DB, token string) (entity.RefreshToken, error)
	DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error
	DeleteByToken(ctx context.Context, tx *gorm.DB, token string) error
	DeleteExpired(ctx context.Context, tx *gorm.DB) error
	FindByPlainToken(ctx context.Context, tx *gorm.DB, token string) (entity.RefreshToken, error)
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) Create(
	ctx context.Context,
	tx *gorm.DB,
	token entity.RefreshToken,
) (entity.RefreshToken, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&token).Error; err != nil {
		return entity.RefreshToken{}, err
	}

	return token, nil
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, tx *gorm.DB, token string) (
	entity.RefreshToken,
	error,
) {
	if tx == nil {
		tx = r.db
	}

	var refreshToken entity.RefreshToken
	if err := tx.WithContext(ctx).Where("token = ?", token).Preload("User").Take(&refreshToken).Error; err != nil {
		return entity.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *refreshTokenRepository) DeleteByUserID(ctx context.Context, tx *gorm.DB, userID string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) DeleteByToken(ctx context.Context, tx *gorm.DB, token string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("token = ?", token).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) DeleteExpired(ctx context.Context, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&entity.RefreshToken{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *refreshTokenRepository) FindByPlainToken(ctx context.Context, tx *gorm.DB, plainToken string) (entity.RefreshToken, error) {
	var tokens []entity.RefreshToken
	if err := tx.WithContext(ctx).Preload("User").Find(&tokens).Error; err != nil {
		return entity.RefreshToken{}, err
	}

	for _, t := range tokens {
		if ok, _ := helpers.CheckPassword(t.Token, []byte(plainToken)); ok {
			return t, nil
		}
	}

	return entity.RefreshToken{}, gorm.ErrRecordNotFound
}
