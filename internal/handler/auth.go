package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "username dan password wajib diisi", err)
		return
	}

	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	response.OK(c, "login berhasil", gin.H{"token": token})
}
