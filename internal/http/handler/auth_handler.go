package handler

import (
	"net/http"
	"shortivo/internal/app/auth"
	appHttp "shortivo/internal/http"
	"shortivo/internal/http/request"
	"shortivo/internal/http/response"
	"strings"
	"time"

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
	token, ok := refreshTokenFromRequest(c)
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

func (h *AuthHandler) SignOut(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.Error(response.BadRequestError("userId is required"))
		return
	}

	h.authService.SignOut(c.Request.Context(), userID)

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

func refreshTokenFromRequest(c *gin.Context) (string, bool) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err == nil && strings.TrimSpace(req.RefreshToken) != "" {
		return req.RefreshToken, true
	}

	authz := strings.TrimSpace(c.GetHeader("Authorization"))
	if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
		token := strings.TrimSpace(authz[len("Bearer "):])
		if token != "" {
			return token, true
		}
	}

	token := strings.TrimSpace(c.GetHeader("X-Refresh-Token"))
	if token != "" {
		return token, true
	}

	return "", false
}
