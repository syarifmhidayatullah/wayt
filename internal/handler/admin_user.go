package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/model"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type AdminUserHandler struct {
	svc service.AuthService
}

func NewAdminUserHandler(svc service.AuthService) *AdminUserHandler {
	return &AdminUserHandler{svc: svc}
}

type createUserRequest struct {
	Username string          `json:"username" binding:"required"`
	Password string          `json:"password" binding:"required"`
	Role     model.AdminRole `json:"role"`
}

type updateUserRequest struct {
	Username string          `json:"username"`
	Password string          `json:"password"`
	Role     model.AdminRole `json:"role"`
}

func (h *AdminUserHandler) List(c *gin.Context) {
	users, err := h.svc.ListUsers()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "success", users)
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	user, err := h.svc.CreateUser(req.Username, req.Password, req.Role)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "user created", user)
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	user, err := h.svc.UpdateUser(uint(id), req.Username, req.Role, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "user updated", user)
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	// ambil requester ID dari JWT claims
	rawID, _ := c.Get("user_id")
	var requesterID uint
	if v, ok := rawID.(float64); ok {
		requesterID = uint(v)
	}

	if err := h.svc.DeleteUser(uint(id), requesterID); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "user deleted", nil)
}
