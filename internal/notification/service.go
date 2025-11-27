package notification

import (
	"context"
	"fmiis/internal/common"
	"fmiis/internal/user"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	FetchCurrentUser(ctx context.Context, email string) (*user.User, error)
	SendSystemEmail(senderEmail, receiverEmail, content, notiType string) error
}

type notificationService struct {
	repo      NotificationRepository
	utilities common.Utilities
}

func NewNotificationService(repo NotificationRepository, utilities common.Utilities) NotificationService {
	return &notificationService{
		repo:      repo,
		utilities: utilities,
	}
}

func (s *notificationService) FetchCurrentUser(ctx context.Context, email string) (*user.User, error) {
	sender, err := s.utilities.GetCurrentUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return sender, nil
}

func (s *notificationService) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	return nil, nil
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

func (s *notificationService) SendSystemEmail(senderEmail, receiverEmail, content, notiType string) error {
	return s.utilities.SendEmail(senderEmail, receiverEmail, content, notiType)
}
