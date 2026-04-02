package handler

import (
	"net/http"
	"strings"
	"time"

	"shortivo/internal/app/auth"
	appHttp "shortivo/internal/http"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cookieManager *appHttp.CookieManager
	authService   auth.AuthService
}

func NewAuthHandler(authService auth.AuthService, cookieManager *appHttp.CookieManager) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		cookieManager: cookieManager,
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.BadRequestError(err.Error()))
		return
	}

	tokens, err := h.authService.SignIn(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.SetRefreshToken(tokens.RefreshToken, tokens.RefreshExp))

	c.JSON(http.StatusOK, response.SuccessResponse[response.SignInResponse]{
		Status:    http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().UTC(),
		Data:      toSignInResponse(tokens),
	})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req request.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(response.BadRequestError("error signup request. Please check request field"))
		return
	}

	tokens, err := h.authService.SignUp(c.Request.Context(), req, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.SetRefreshToken(tokens.RefreshToken, tokens.RefreshExp))

	c.JSON(http.StatusOK, response.SuccessResponse[response.SignInResponse]{
		Status:    http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().UTC(),
		Data:      toSignInResponse(tokens),
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	token, ok := getRefreshToken(c)
	if !ok {
		c.Error(response.BadRequestError("refresh token is required"))
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), token, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.SetRefreshToken(tokens.RefreshToken, tokens.RefreshExp))

	c.JSON(http.StatusOK, response.SuccessResponse[response.SignInResponse]{
		Status:    http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().UTC(),
		Data:      toSignInResponse(tokens),
	})
}

func (h *AuthHandler) GoogleSignIn(c *gin.Context) {
	loginURL, state, err := h.authService.GoogleSignInURL()
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.SetGoogleOAuthState(state))
	c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	req := request.GoogleOAuthCallbackRequest{
		Code:             strings.TrimSpace(c.Query("code")),
		State:            strings.TrimSpace(c.Query("state")),
		Error:            strings.TrimSpace(c.Query("error")),
		ErrorDescription: strings.TrimSpace(c.Query("error_description")),
	}

	if strings.TrimSpace(req.Error) != "" {
		errMsg := req.ErrorDescription
		if errMsg == "" {
			errMsg = req.Error
		}
		c.Error(response.BadRequestError(errMsg))
		return
	}

	if req.Code == "" || req.State == "" {
		c.Error(response.BadRequestError("invalid google oauth callback"))
		return
	}

	state, err := c.Cookie(appHttp.GoogleOAuthStateCookieName())
	if err != nil || strings.TrimSpace(state) == "" {
		c.Error(response.Unauthorized())
		return
	}

	c.SetCookieData(h.cookieManager.RevokeGoogleOAuthState())

	tokens, err := h.authService.SignInWithGoogle(c.Request.Context(), req, state, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.SetRefreshToken(tokens.RefreshToken, tokens.RefreshExp))

	c.JSON(http.StatusOK, response.SuccessResponse[response.SignInResponse]{
		Status:    http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().UTC(),
		Data:      toSignInResponse(tokens),
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.Error(response.BadRequestError("userId is required"))
		return
	}

	token, _ := getRefreshToken(c)

	if err := h.authService.SignOut(c.Request.Context(), userID, token); err != nil {
		c.Error(err)
		return
	}

	c.SetCookieData(h.cookieManager.RevokeRefreshToken())

	c.JSON(http.StatusOK, response.SuccessResponse[bool]{
		Status:    http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().UTC(),
		Data:      true,
	})
}

func toSignInResponse(tokens auth.AuthTokens) response.SignInResponse {
	return response.SignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func getRefreshToken(c *gin.Context) (string, bool) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil && strings.TrimSpace(req.RefreshToken) != "" {
		return req.RefreshToken, true
	}

	token := strings.TrimSpace(c.GetHeader("X-Refresh-Token"))
	if token != "" {
		return token, true
	}

	token, _ = c.Cookie("refreshToken")
	if token != "" {
		return token, true
	}

	return "", false
}
