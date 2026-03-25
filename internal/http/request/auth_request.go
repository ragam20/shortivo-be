package request

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=5"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SignOutRequest struct {
	UserID int64 `json:"userId" binding:"required"`
}

type GoogleOAuthCallbackRequest struct {
	Code             string `form:"code" binding:"required"`
	State            string `form:"state" binding:"required"`
	Error            string `form:"error"`
	ErrorDescription string `form:"error_description"`
}
