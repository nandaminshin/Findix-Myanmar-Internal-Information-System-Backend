package notification

import (
	"context"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error)
}

type notificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) NotificationService {
	return &notificationService{
		repo: repo,
	}
}

func (s *notificationService) SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error) {
	notification := &Notification{
		NotiType: req.NotiType,
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Content:  req.Content,
		IsSeen:   false,
	}
	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return &NotificationResponse{
		NotiType:  notification.NotiType,
		Sender:    notification.Sender,
		Receiver:  notification.Receiver,
		Content:   notification.Content,
		IsSeen:    notification.IsSeen,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}, nil
}
