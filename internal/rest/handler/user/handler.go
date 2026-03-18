package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry/internal/config"
	"github.com/robiuzzaman4/donor-registry/internal/domain"
	"github.com/robiuzzaman4/donor-registry/internal/rest/middleware"
	"github.com/robiuzzaman4/donor-registry/internal/rest/response"
	"github.com/robiuzzaman4/donor-registry/internal/util"
)

type Handler struct {
	cnf *config.Config
	svc Service
}

func NewHandler(cnf *config.Config, svc Service) *Handler {
	return &Handler{
		cnf: cnf,
		svc: svc,
	}
}

/*
================================================
Auth Handlers
================================================
*/

// Register handler / Create user handler
func (h *Handler) Register(c *gin.Context) {

	var req struct {
		Name             string                `json:"name" binding:"required"`
		Phone            string                `json:"phone" binding:"required"`
		Password         string                `json:"password" binding:"required"`
		BloodGroup       domain.UserBloodGroup `json:"blood_group" binding:"required"`
		Role             domain.UserRole       `json:"role"`
		Gender           domain.UserGender     `json:"gender" binding:"required"`
		DateOfBirth      time.Time             `json:"date_of_birth" binding:"required"`
		Zila             string                `json:"zila"`
		Upazila          string                `json:"upazila"`
		LocalAddress     string                `json:"local_address"`
		TotalDonateCount int                   `json:"total_donate_count"`
		IsAvailable      bool                  `json:"is_available"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		response.BadRequest(c, "Failed to process password")
		return
	}
	user := domain.User{
		Name:             req.Name,
		Phone:            req.Phone,
		Password:         hashedPassword,
		BloodGroup:       req.BloodGroup,
		Role:             req.Role,
		Gender:           req.Gender,
		DateOfBirth:      req.DateOfBirth,
		Zila:             req.Zila,
		Upazila:          req.Upazila,
		LocalAddress:     req.LocalAddress,
		TotalDonateCount: req.TotalDonateCount,
		IsAvailable:      req.IsAvailable,
	}

	res, err := h.svc.Create(c.Request.Context(), &user)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, res)
}

// Login handler
func (h *Handler) Login(c *gin.Context) {
	var loginReq struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	type loginRes struct {
		AccessToken  string       `json:"access_token"`
		RefreshToken string       `json:"refresh_token"`
		User         *domain.User `json:"user"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.BadRequest(c, "Phone and password are required")
		return
	}

	loginReq.Phone = strings.TrimSpace(loginReq.Phone)
	loginReq.Password = strings.TrimSpace(loginReq.Password)

	user, err := h.svc.GetByPhone(c.Request.Context(), loginReq.Phone)
	if err != nil {
		response.Error(c, err)
		return
	}

	if user == nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	fmt.Println("user.Password", user.Password)
	fmt.Println("loginReq.Password", loginReq.Password)

	match := util.CheckPasswordHash(loginReq.Password, user.Password)
	fmt.Println("match", match)
	if !match {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	accessTokenExpiry := 24 * time.Hour
	refreshTokenExpiry := 168 * time.Hour

	accessToken, err := util.GenerateToken(user.ID, string(user.Role), accessTokenExpiry)
	if err != nil {
		response.Error(c, err)
		return
	}
	refreshToken, err := util.GenerateToken(user.ID, string(user.Role), refreshTokenExpiry)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.SetCookie(
		"access_token",                   // name
		accessToken,                      // value
		int(accessTokenExpiry.Seconds()), // maxAge in seconds
		"/",                              // path
		"",                               // domain (empty for current domain)
		false,                            // secure (set to true in production/HTTPS)
		true,                             // httpOnly (prevents JS access)
	)
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(refreshTokenExpiry.Seconds()),
		"/",
		"",
		false, // secure
		true,  // httpOnly
	)

	response.SuccessWithMessage(c, "Login successful", loginRes{
		accessToken,
		refreshToken,
		user,
	})
}

// Refresh access token using refresh_token handler
func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	type refreshRes struct {
		AccessToken string `json:"access_token"`
	}

	_ = c.ShouldBindJSON(&req)
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		if cookieToken, err := c.Cookie("refresh_token"); err == nil {
			refreshToken = strings.TrimSpace(cookieToken)
		}
	}
	if refreshToken == "" {
		response.BadRequest(c, "Refresh token is required")
		return
	}

	claims, err := util.ValidateToken(c.Request.Context(), refreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	userID := claims.ID
	if userID == "" {
		response.Error(c, domain.ErrInvalidToken)
		return
	}
	role := claims.Role

	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	if user == nil {
		response.Error(c, domain.ErrUserNotFound)
		return
	}

	accessTokenExpiry := 24 * time.Hour
	accessToken, err := util.GenerateToken(userID, role, accessTokenExpiry)
	if err != nil {
		response.Error(c, err)
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		int(accessTokenExpiry.Seconds()),
		"/",
		"",
		false,
		true,
	)

	response.SuccessWithMessage(c, "Access token refreshed", refreshRes{
		AccessToken: accessToken,
	})
}

// Logout handler
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	response.SuccessWithMessage(c, "Logout successful", nil)
}

/*
================================================
User Handlers - Public
================================================
*/

// List users handler
func (h *Handler) List(c *gin.Context) {
	page, limit, _ := response.ParsePaginationParams(c)

	users, total, err := h.svc.List(c.Request.Context(), page, limit)
	if err != nil {
		response.Error(c, err)
		return
	}

	pagination := response.BuildPagination(total, page, limit)
	response.SuccessWithPagination(c, "Donors retrieved successfully", users, pagination)
}

/*
================================================
User Handlers - Protectd
================================================
*/

// Get by id handler
func (h *Handler) GetByID(c *gin.Context) {
	userID, isFound := c.Params.Get("userID")
	if userID == "" || !isFound {
		response.BadRequest(c, "UserID is required")
		return
	}

	user, err := h.svc.GetByID(c, userID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessWithMessage(c, "User retrived", user)
}

// Get by phone handler
func (h *Handler) GetByPhone(c *gin.Context) {
	phone, isFound := c.Params.Get("phone")
	if phone == "" || !isFound {
		response.BadRequest(c, "Phone is required")
		return
	}

	user, err := h.svc.GetByPhone(c, phone)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessWithMessage(c, "User retrived", user)
}

// Get current user handler
func (h *Handler) Me(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		response.Error(c, domain.ErrUnauthorized)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.Error(c, domain.ErrInvalidToken)
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), userIDStr)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.SuccessWithMessage(c, "User retrived", user)
}

// Update user handler
func (h *Handler) Update(c *gin.Context) {
	userID, isFound := c.Params.Get("userID")
	if userID == "" || !isFound {
		response.BadRequest(c, "UserID is required")
		return
	}

	var req domain.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	err := h.svc.Update(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "User Updated", &req)
}
