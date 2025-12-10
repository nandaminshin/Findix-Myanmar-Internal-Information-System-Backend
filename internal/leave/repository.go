package leave

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LeaveRepository interface {
	FindBymEmpIDAndEndDate(ctx context.Context, empID primitive.ObjectID, startDate primitive.DateTime, endDate primitive.DateTime) (*Leave, error)
	Create(ctx context.Context, leave *Leave) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*Leave, error)
	Update(ctx context.Context, leave *Leave) error
}

type leaveRepository struct {
	collection *mongo.Collection
}

func NewLeaveRepository(db *mongo.Database) LeaveRepository {
	return &leaveRepository{
		collection: db.Collection("leaves"),
	}
}

func (r *leaveRepository) Create(ctx context.Context, leave *Leave) error {
	leave.ID = primitive.NewObjectID()
	leave.CreatedAt = time.Now()
	leave.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, &leave)
	return err
}

func (r *leaveRepository) FindBymEmpIDAndEndDate(ctx context.Context, empID primitive.ObjectID, startDate primitive.DateTime, endDate primitive.DateTime) (*Leave, error) {
	var leave Leave

	filter := bson.M{
		"empID":     empID,
		"status":    "confirmed",
		"startDate": bson.M{"$lte": startDate},
		"$and": []bson.M{
			{"startDate": bson.M{"$lte": endDate}},
			{"endDate": bson.M{"$gte": startDate}},
		},
	}

	err := r.collection.FindOne(ctx, filter).Decode(&leave)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &leave, nil
}

func (r *leaveRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Leave, error) {
	var leave Leave
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&leave)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &leave, nil
}

func (r *leaveRepository) Update(ctx context.Context, leave *Leave) error {
	_, err := r.collection.UpdateByID(ctx, &leave.ID,
		bson.M{
			"$set": &leave,
		},
	)
	return err
}
