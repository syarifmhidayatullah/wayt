package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/project/wayt/internal/service"
	"github.com/project/wayt/pkg/response"
)

type UserHandler struct {
	userSvc  service.UserService
	queueSvc service.QueueService
}

func NewUserHandler(userSvc service.UserService, queueSvc service.QueueService) *UserHandler {
	return &UserHandler{userSvc: userSvc, queueSvc: queueSvc}
}

type userRegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userLoginRequest struct {
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req userRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	user, err := h.userSvc.Register(req.Name, req.Phone, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "registrasi berhasil", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"phone": user.Phone,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req userLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request", err)
		return
	}
	token, err := h.userSvc.Login(req.Phone, req.Password)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.OK(c, "login berhasil", gin.H{"token": token})
}

// ListBranches — public: semua branch aktif beserta counter & jumlah antrian
func (h *UserHandler) ListBranches(c *gin.Context) {
	search := c.Query("search")
	branches, err := h.queueSvc.ListPublicBranches(search)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "success", branches)
}

// BookQueue — authenticated user books a queue at a counter
func (h *UserHandler) BookQueue(c *gin.Context) {
	rawID, _ := c.Get("user_id")
	var userID uint
	if v, ok := rawID.(float64); ok {
		userID = uint(v)
	}

	var req struct {
		CounterID uint `json:"counter_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "counter_id wajib diisi", err)
		return
	}

	result, err := h.queueSvc.BookByUser(userID, req.CounterID)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Created(c, "antrian berhasil diambil", result)
}

// MyQueues — authenticated user's active queues
func (h *UserHandler) MyQueues(c *gin.Context) {
	rawID, _ := c.Get("user_id")
	var userID uint
	if v, ok := rawID.(float64); ok {
		userID = uint(v)
	}
	queues, err := h.queueSvc.MyQueues(userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "success", queues)
}

// DashboardPage renders the user-facing dashboard HTML
func (h *UserHandler) DashboardPage(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{})
}
