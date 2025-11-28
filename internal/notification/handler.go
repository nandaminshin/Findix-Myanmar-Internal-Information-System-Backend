package notification

import (
	"fmt"
	"log"
	"net/http"

	"fmiis/internal/user"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service     NotificationService
	userService user.UserService
}

func NewNotificationHandler(service NotificationService, userService user.UserService) *NotificationHandler {
	return &NotificationHandler{
		service:     service,
		userService: userService,
	}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Receivers:", req.Receivers)

	// fetch sender by sender email
	sender, err := h.service.FetchCurrentUser(c.Request.Context(), req.Sender.Email)
	if err != nil {
		log.Printf("Sender not found with email %s: %v", req.Sender.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender not found"})
		return
	}

	receivers, err := h.service.FetchReceivers(c.Request.Context(), req.Receivers)
	if err != nil {
		log.Printf("Error fetching receivers: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error fetching receivers"})
		return
	}

	res, err := h.service.SendNotification(c.Request.Context(), &NotificationRequest{
		NotiType:  req.NotiType,
		Sender:    *sender,
		Receivers: receivers,
		Content:   req.Content,
		IsSeen:    req.IsSeen,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SendEmailNotificationForAllReceivers(receivers, res)
	if err != nil {
		log.Printf("Error sending email notifications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending email notifications"})
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}
