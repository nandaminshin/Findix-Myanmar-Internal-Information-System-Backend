package leave

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaveHandler struct {
	service LeaveService
}

func NewLeaveHandler(service LeaveService) *LeaveHandler {
	return &LeaveHandler{
		service: service,
	}
}

func (h *LeaveHandler) CreateLeaveRequest(c *gin.Context) {
	var req LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.RecordLeaveRequest(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *LeaveHandler) LeaveRequestGmApproval(c *gin.Context) {
	var req GmApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.GmApproval(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
