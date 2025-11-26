package notification

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	// FindByReceiver(ctx context.Context, receiverID primitive.ObjectID) ([]*Notification, error)
	// FindBySender(ctx context.Context, senderID primitive.ObjectID) ([]*Notification, error)
}

type mongoNotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &mongoNotificationRepository{
		collection: db.Collection("notifications"),
	}
}

func (r *mongoNotificationRepository) Create(ctx context.Context, notification *Notification) error {
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}
