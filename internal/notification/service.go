package notification

import (
	"context"
	"fmiis/internal/common"
	"fmiis/internal/normal_email"
	"fmiis/internal/user"
	"fmt"
	"html"
	"log"
	"sync"
	"time"
)

type NotificationService interface {
	SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error)
	SendEmailNotificationForAllReceivers(receiver []Receiver, notification *NotificationResponse) error
	FetchCurrentUser(ctx context.Context, email string) (*user.User, error)
	FetchReceivers(ctx context.Context, receivers []Receiver) ([]Receiver, error)
	FormatNotificationHTML(notification NotificationResponse) (string, string)
	formatSenderPhone(phone string) string
	formatEditedTime(createdAt, updatedAt time.Time) string
}

type notificationService struct {
	repo         NotificationRepository
	utilities    common.Utilities
	emailService normal_email.EmailService
	userService  user.UserService
}

func NewNotificationService(repo NotificationRepository, utilities common.Utilities, emailService normal_email.EmailService, userService user.UserService) NotificationService {
	return &notificationService{
		repo:         repo,
		utilities:    utilities,
		emailService: emailService,
		userService:  userService,
	}
}

func (s *notificationService) FetchCurrentUser(ctx context.Context, email string) (*user.User, error) {
	sender, err := s.utilities.GetCurrentUser(ctx, email)
	if err != nil {
		return nil, err
	}
	return sender, nil
}

func (s *notificationService) FetchReceivers(ctx context.Context, receivers []Receiver) ([]Receiver, error) {
	// fetch receiver by receiver email
	var wg sync.WaitGroup
	var mu sync.Mutex
	var receiverErr error
	var fetchedReceivers []Receiver

	for i := range receivers { // Use index
		wg.Add(1)
		go func(idx int) { // Pass index
			defer wg.Done()

			r := receivers[idx]
			receiverUser, err := s.userService.GetUserByEmail(ctx, r.Email)
			if err != nil {
				mu.Lock()
				if receiverErr == nil {
					receiverErr = fmt.Errorf("receiver not found: %s", r.Email)
				}
				mu.Unlock()
				return
			}

			if receiverUser == nil {
				mu.Lock()
				if receiverErr == nil {
					receiverErr = fmt.Errorf("receiver not found: %s", r.Email)
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			fetchedReceivers = append(fetchedReceivers, Receiver{
				ID:    receiverUser.ID,
				Name:  receiverUser.Name,
				Email: receiverUser.Email,
			})
			mu.Unlock()
		}(i) // Pass current index
	}
	wg.Wait()
	if receiverErr != nil {
		return nil, receiverErr
	}
	return receivers, nil
}

func (s *notificationService) SendNotification(ctx context.Context, req *NotificationRequest) (*NotificationResponse, error) {
	notification := &Notification{
		NotiType:  req.NotiType,
		Sender:    req.Sender,
		Receivers: req.Receivers,
		Content:   req.Content,
		IsSeen:    false,
	}

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return &NotificationResponse{
		NotiType:  notification.NotiType,
		Sender:    notification.Sender,
		Receivers: notification.Receivers,
		Content:   notification.Content,
		IsSeen:    notification.IsSeen,
		CreatedAt: notification.CreatedAt,
		UpdatedAt: notification.UpdatedAt,
	}, nil
}

func (s *notificationService) SendEmailNotificationForAllReceivers(receivers []Receiver, notification *NotificationResponse) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var emailErr error

	for i := range receivers { // Use index instead of value
		wg.Add(1)
		go func(idx int) { // Pass index to goroutine
			defer wg.Done()

			r := receivers[idx] // Get receiver inside goroutine

			err := s.SendEmailNotification(r.Email, r.Name, notification)
			if err != nil {
				mu.Lock()
				if emailErr == nil {
					emailErr = fmt.Errorf("failed to send to %s: %v", r.Email, err)
				} else {
					emailErr = fmt.Errorf("%w; failed to send to %s: %v", emailErr, r.Email, err)
				}
				mu.Unlock()
			}
		}(i) // Pass the current index
	}
	wg.Wait()
	if emailErr != nil {
		return emailErr
	}
	return nil
}

// Email seiding things
func (s *notificationService) FormatNotificationHTML(notification NotificationResponse) (string, string) {
	subject := fmt.Sprintf("🔔 %s  - FMIIS", notification.NotiType)

	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 0; background: #f5f5f5; }
        .container { background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin: 20px auto; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; }
        .content { padding: 30px; }
        .notification-type { font-size: 24px; font-weight: bold; margin-bottom: 10px; color: #2c3e50; }
        .message-box { background: #f8f9fa; padding: 20px; border-radius: 8px; border-left: 4px solid #667eea; margin: 20px 0; }
        .info-box { background: white; padding: 20px; border-radius: 8px; border: 1px solid #e9ecef; margin: 20px 0; }
        .timestamp { color: #6c757d; font-size: 14px; margin-top: 20px; padding-top: 20px; border-top: 1px solid #e9ecef; }
        .badge { display: inline-block; background: #667eea; color: white; padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: 500; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #e9ecef; color: #6c757d; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 style="margin: 0; font-size: 28px;">🔔 Findix Myanmar Internal Information System</h1>
            <p style="margin: 5px 0 0 0; opacity: 0.9;">Email Notification System</p>
        </div>
        
        <div class="content">
            <div class="notification-type">%s Notification</div>
            
            <div class="message-box">
                <h3 style="margin-top: 0; color: #2c3e50;">📝 Message</h3>
                <p style="margin: 0; font-size: 16px; line-height: 1.5;">%s</p>
            </div>

            <div class="info-box">
                <h3 style="margin-top: 0; color: #2c3e50;">👤 From</h3>
                <table style="width: 100%%; border-collapse: collapse;">
                    <tr><td style="padding: 8px 0; width: 80px;"><strong>Name:</strong></td><td style="padding: 8px 0;">%s</td></tr>
                    <tr><td style="padding: 8px 0;"><strong>Role:</strong></td><td style="padding: 8px 0;"><span class="badge">%s</span></td></tr>
                    <tr><td style="padding: 8px 0;"><strong>Email:</strong></td><td style="padding: 8px 0;">%s</td></tr>
                    %s
                </table>
            </div>

            <div class="timestamp">
                <p style="margin: 5px 0;"><strong>🕒 Sent:</strong> %s</p>
                %s
            </div>
        </div>
        
        <div class="footer">
            <p style="margin: 0;">This is an automated notification from Findix FMIIS System</p>
            <p style="margin: 5px 0;">© %d Findix Myanmar. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
		// Notification type
		notification.NotiType,

		// Message content (escaped for HTML safety)
		html.EscapeString(notification.Content),

		// Sender information
		html.EscapeString(notification.Sender.Name),
		html.EscapeString(string(notification.Sender.Role)),
		html.EscapeString(notification.Sender.Email),

		// Sender phone (conditional)
		s.formatSenderPhone(notification.Sender.Phone),

		// Timestamps
		notification.CreatedAt.Format("January 2, 2006 at 3:04 PM"),

		// Edited timestamp (conditional)
		s.formatEditedTime(notification.CreatedAt, notification.UpdatedAt),

		// Copyright year
		time.Now().Year(),
	)

	return subject, htmlContent
}

func (s *notificationService) formatSenderPhone(phone string) string {
	if phone == "" {
		return ""
	}
	return fmt.Sprintf("<tr><td style='padding: 8px 0;'><strong>Phone:</strong></td><td style='padding: 8px 0;'>%s</td></tr>", html.EscapeString(phone))
}

func (s *notificationService) formatEditedTime(createdAt, updatedAt time.Time) string {
	if !updatedAt.Equal(createdAt) && updatedAt.After(createdAt) {
		return fmt.Sprintf("<p style='margin: 5px 0;'><strong>✏️ Last Edited:</strong> %s</p>", updatedAt.Format("January 2, 2006 at 3:04 PM"))
	}
	return ""
}

func (s *notificationService) SendEmailNotification(receiverEmail, receiverName string, notification *NotificationResponse) error {
	if receiverEmail == "" {
		log.Printf("No email address for receiver")
		return fmt.Errorf("no email address for receiver")
	}

	// Format email
	subject, htmlContent := s.FormatNotificationHTML(*notification)

	// Send email
	err := s.emailService.SendNotificationEmail(
		receiverEmail,
		receiverName,
		subject,
		htmlContent,
	)

	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", receiverEmail, err)
		return err
	} else {
		log.Printf("✅ Email sent successfully to %s", receiverEmail)
		return nil
	}
}
