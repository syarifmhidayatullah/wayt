package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type QueueHandler struct {
	svc service.QueueService
}

func NewQueueHandler(svc service.QueueService) *QueueHandler {
	return &QueueHandler{svc: svc}
}

type registerRequest struct {
	QRToken string `json:"qr_token" binding:"required"`
}

// ScanRegister dipanggil langsung saat QR di-scan — register lalu redirect ke halaman status per queue ID
func (h *QueueHandler) ScanRegister(c *gin.Context) {
	token := c.Param("token")
	result, err := h.svc.ScanRegister(token)
	if err != nil {
		c.HTML(200, "queue.html", gin.H{"error": err.Error()})
		return
	}
	c.Redirect(302, fmt.Sprintf("/queue/%d", result.QueueID))
}

// QueuePage render halaman status antrian berdasarkan queue ID
func (h *QueueHandler) QueuePage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(200, "queue.html", gin.H{"error": "link tidak valid"})
		return
	}
	result, err := h.svc.StatusByID(uint(id))
	if err != nil {
		c.HTML(200, "queue.html", gin.H{"error": err.Error()})
		return
	}
	c.HTML(200, "queue.html", gin.H{
		"queue_id":     id,
		"queue_number": result.QueueNumber,
		"branch_name":  "",
		"status":       string(result.Status),
		"current":      result.CurrentServing,
		"people_ahead": result.PeopleAhead,
	})
}

func (h *QueueHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	result, err := h.svc.Register(req.QRToken)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "queue registered", result)
}

func (h *QueueHandler) Status(c *gin.Context) {
	token := c.Param("token")
	result, err := h.svc.Status(token)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "success", result)
}

func (h *QueueHandler) StatusByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	result, err := h.svc.StatusByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "success", result)
}
