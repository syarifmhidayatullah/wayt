package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type QRCodeHandler struct {
	svc service.QRCodeService
}

func NewQRCodeHandler(svc service.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{svc: svc}
}

func (h *QRCodeHandler) Generate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid branch id", err)
		return
	}
	result, err := h.svc.Generate(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "QR code generated", result)
}
