package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type BranchHandler struct {
	svc service.BranchService
}

func NewBranchHandler(svc service.BranchService) *BranchHandler {
	return &BranchHandler{svc: svc}
}

type createBranchRequest struct {
	Name string `json:"name" binding:"required"`
}

type updateBranchRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

func (h *BranchHandler) Create(c *gin.Context) {
	var req createBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	branch, err := h.svc.Create(req.Name)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "branch created", branch)
}

// List returns branches filtered by role:
// - superadmin: all branches
// - admin: only their assigned branch
func (h *BranchHandler) List(c *gin.Context) {
	role, _ := c.Get("role")

	if role == "admin" {
		rawBranchID, exists := c.Get("branch_id")
		if !exists || rawBranchID == nil {
			response.OK(c, "success", []interface{}{})
			return
		}
		var branchID uint
		switch v := rawBranchID.(type) {
		case float64:
			branchID = uint(v)
		case uint:
			branchID = v
		}
		if branchID == 0 {
			response.OK(c, "success", []interface{}{})
			return
		}
		branch, err := h.svc.FindByID(branchID)
		if err != nil {
			response.OK(c, "success", []interface{}{})
			return
		}
		response.OK(c, "success", []interface{}{branch})
		return
	}

	branches, err := h.svc.List()
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "success", branches)
}

func (h *BranchHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	var req updateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	branch, err := h.svc.Update(uint(id), req.Name, isActive)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "branch updated", branch)
}

func (h *BranchHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "branch deleted", nil)
}
