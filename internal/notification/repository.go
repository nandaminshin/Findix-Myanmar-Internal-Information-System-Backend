package notification

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	SetupTTLIndex(ctx context.Context) error
	FindAllByReceiver(ctx context.Context, receiverEmail string) ([]Notification, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Notification, error)
	UpdateByID(ctx context.Context, id primitive.ObjectID, notification *Notification) error
	// FindBySender(ctx context.Context, senderID primitive.ObjectID) ([]*Notification, error)
}

type mongoNotificationRepository struct {
	collection *mongo.Collection
}

func IsValidNotiType(notiType NotiType) bool {
	switch notiType {
	case morningMeetingNnoti, devMeetingNoti, kosugiMeeting, emergencyMeetingNoti, internalMeetingNoti, generalNoti:
		return true
	default:
		return false
	}
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &mongoNotificationRepository{
		collection: db.Collection("notifications"),
	}
}

// SetupTTLIndex creates an index that automatically deletes notifications after 7 days
func (r *mongoNotificationRepository) SetupTTLIndex(ctx context.Context) error {
	// 1 days in seconds (1 * 24 * 60 * 60)
	expireAfterSeconds := int32(7 * 24 * 60 * 60)

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
	if !IsValidNotiType(notification.NotiType) {
		return errors.New("invalid notification type")
	}
	notification.ID = primitive.NewObjectID()
	notification.CreatedAt = time.Now()
	notification.UpdatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, &notification)
	return err
}

func (r *mongoNotificationRepository) FindAllByReceiver(ctx context.Context, receiverEmail string) ([]Notification, error) {
	opt := options.Find().SetSort(bson.M{
		"created_at": -1,
	})
	cursor, err := r.collection.Find(ctx, bson.M{
		"receivers.email": receiverEmail,
	}, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *mongoNotificationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Notification, error) {
	res := &Notification{}
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *mongoNotificationRepository) UpdateByID(ctx context.Context, id primitive.ObjectID, notification *Notification) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{
		"_id": id,
	}, bson.M{
		"$set": notification,
	})
	return err
}
