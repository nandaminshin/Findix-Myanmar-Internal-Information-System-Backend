package notification

import (
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

	// fetch sender by sender email
	sender, err := h.service.FetchCurrentUser(c.Request.Context(), req.Sender.Email)
	if err != nil {
		log.Printf("Sender not found with email %s: %v", req.Sender.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender not found"})
		return
	}

	// fetch receiver by receiver email
	receiver, err := h.userService.GetUserByEmail(c.Request.Context(), req.Receiver.Email)
	if err != nil {
		log.Printf("Receiver not found with email %s: %v", req.Receiver.Email, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver not found"})
		return
	}

	res, err := h.service.SendNotification(c.Request.Context(), &NotificationRequest{
		NotiType: req.NotiType,
		Sender:   *sender,
		Receiver: *receiver,
		Content:  req.Content,
		IsSeen:   req.IsSeen,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.SendSystemEmail(sender.Email, receiver.Email, req.Content, string(req.NotiType))
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
	}

	c.IndentedJSON(http.StatusOK, res)
}
