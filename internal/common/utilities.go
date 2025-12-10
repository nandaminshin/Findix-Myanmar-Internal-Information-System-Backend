package common

import (
	"context"
	"errors"
	"fmiis/internal/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (u *Utilities) ParseDateTimeForDB(dt string) (*primitive.DateTime, error) {
	t, err := time.Parse("02/01/2006", dt) //dd/mm/yyyy format date
	if err != nil {
		return nil, errors.New("invalid date format, expected dd/mm/yy")
	}
	// Convert to primitive.DateTime for MongoDB
	dbDateTime := primitive.NewDateTimeFromTime(t)
	return &dbDateTime, nil
}

func (u *Utilities) ParseObjectIDForDB(id string) (*primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	return &objID, nil
}
