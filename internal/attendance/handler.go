package attendance

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AttendanceHandler struct {
	service AttendanceService
}

func NewAttendanceHandler(service AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		service: service,
	}
}

func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	var req AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.RecordAttendance(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
