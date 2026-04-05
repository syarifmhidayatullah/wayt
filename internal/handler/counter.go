package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type CounterHandler struct {
	counterSvc service.CounterService
	qrSvc      service.QRCodeService
	queueSvc   service.QueueService
}

func NewCounterHandler(counterSvc service.CounterService, qrSvc service.QRCodeService, queueSvc service.QueueService) *CounterHandler {
	return &CounterHandler{counterSvc: counterSvc, qrSvc: qrSvc, queueSvc: queueSvc}
}

type createCounterRequest struct {
	Name   string `json:"name"   binding:"required"`
	Prefix string `json:"prefix" binding:"required"`
}

type updateCounterRequest struct {
	Name     string `json:"name"`
	Prefix   string `json:"prefix"`
	IsActive *bool  `json:"is_active"`
}

// adminBranchID returns the branch_id from JWT context (for admin role only).
// Returns 0 if not set (superadmin or unassigned).
func adminBranchID(c *gin.Context) uint {
	rawBranchID, exists := c.Get("branch_id")
	if !exists || rawBranchID == nil {
		return 0
	}
	switch v := rawBranchID.(type) {
	case float64:
		return uint(v)
	case uint:
		return v
	}
	return 0
}

// checkBranchAccess returns false and writes 403 if the admin user's branch_id
// does not match the requested branchID. Superadmins always pass.
func checkBranchAccess(c *gin.Context, branchID uint) bool {
	role, _ := c.Get("role")
	if role == "superadmin" {
		return true
	}
	if adminBranchID(c) != branchID {
		response.Forbidden(c)
		return false
	}
	return true
}

func (h *CounterHandler) Create(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid branch id", err)
		return
	}
	if !checkBranchAccess(c, uint(branchID)) {
		return
	}
	var req createCounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	counter, err := h.counterSvc.Create(uint(branchID), req.Name, req.Prefix)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "counter created", counter)
}

func (h *CounterHandler) ListByBranch(c *gin.Context) {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid branch id", err)
		return
	}
	if !checkBranchAccess(c, uint(branchID)) {
		return
	}
	counters, err := h.counterSvc.ListByBranch(uint(branchID))
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "success", counters)
}

func (h *CounterHandler) Update(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	// verify access by looking up counter's branch
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	var req updateCounterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	updated, err := h.counterSvc.Update(uint(counterID), req.Name, req.Prefix, isActive)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "counter updated", updated)
}

func (h *CounterHandler) Delete(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	if err := h.counterSvc.Delete(uint(counterID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "counter deleted", nil)
}

func (h *CounterHandler) GenerateQR(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	result, err := h.qrSvc.Generate(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "QR generated", result)
}

func (h *CounterHandler) CallNext(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	queue, err := h.queueSvc.CallNext(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "next queue called", queue)
}

func (h *CounterHandler) ListQueue(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	queues, err := h.queueSvc.ListByCounter(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "success", queues)
}

func (h *CounterHandler) Reset(c *gin.Context) {
	counterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	counter, err := h.counterSvc.FindByID(uint(counterID))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	if !checkBranchAccess(c, counter.BranchID) {
		return
	}
	if err := h.queueSvc.Reset(uint(counterID)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "counter reset successfully", nil)
}
