package user

import (
	"github.com/gin-gonic/gin"
	"github.com/robiuzzaman4/donor-registry/internal/config"
	"github.com/robiuzzaman4/donor-registry/internal/domain"
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

// register handler / create user handler
func (h *Handler) Register(c *gin.Context) {

	var req domain.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		response.BadRequest(c, "Failed to process password")
		return
	}
	req.Password = hashedPassword

	res, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, res)
}

// login user handler
func (h *Handler) Login(c *gin.Context) {
	var loginReq struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		response.BadRequest(c, "Phone and password are required")
		return
	}

	user, err := h.svc.FindByPhoneAndPassword(c.Request.Context(), loginReq.Phone, loginReq.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	if user == nil {
		response.Unauthorized(c, "Invalid phone number or password")
		return
	}

	response.SuccessWithMessage(c, "Login successful", user)
}

// list users
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

// get by id
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

// get by phone
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

// update
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
