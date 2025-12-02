package notification

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	SetupTTLIndex(ctx context.Context) error
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

// SetupTTLIndex creates an index that automatically deletes notifications after 7 days
func (r *mongoNotificationRepository) SetupTTLIndex(ctx context.Context) error {
	// 7 days in seconds (7 * 24 * 60 * 60)
	expireAfterSeconds := int32(2 * 60)

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			"created_at": 1, // Index on the created_at field
		},
		Options: options.Index().
			SetExpireAfterSeconds(expireAfterSeconds). // Auto-delete after 7 days
			SetName("notification_ttl_index"),         // Optional: name the index
	}

	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *mongoNotificationRepository) Create(ctx context.Context, notification *Notification) error {
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, notification)
	return err
}
