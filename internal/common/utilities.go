package common

import (
	"context"
	"fmiis/internal/user"
	"os"

	"github.com/resend/resend-go/v2"
)

type Utilities struct {
	userService user.UserService
}

func NewUtility(service user.UserService) *Utilities {
	return &Utilities{
		userService: service,
	}
}

func (u *Utilities) GetCurrentUser(ctx context.Context, email string) (*user.User, error) {
	user, err := u.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Utilities) SendEmail(senderEmail, receiverEmail, content, notiType string) error {
	apiKey := os.Getenv("RESEND_API")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		To:      []string{receiverEmail},
		From:    senderEmail,
		Text:    content,
		Subject: notiType,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
