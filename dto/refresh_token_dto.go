package dto

const (
	MESSAGE_SUCCESS_REFRESH_TOKEN        = "successfully refreshed token"
	MESSAGE_FAILED_REFRESH_TOKEN         = "failed to refresh token"
	MESSAGE_FAILED_INVALID_REFRESH_TOKEN = "invalid refresh token"
	MESSAGE_FAILED_EXPIRED_REFRESH_TOKEN = "refresh token has expired"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Role         string `json:"role"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
