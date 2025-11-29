package common

import (
	"context"
	"fmiis/internal/user"
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
