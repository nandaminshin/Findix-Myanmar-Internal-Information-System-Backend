package attendance

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AttendanceRepository interface {
	Create(ctx context.Context, attendance *Attendance) error
	FindBymEmpIDAndDate(ctx context.Context, empID primitive.ObjectID, date primitive.DateTime) (*Attendance, error)
}

type attendanceRepository struct {
	collection *mongo.Collection
}

func NewAttendanceRepository(db *mongo.Database) AttendanceRepository {
	return &attendanceRepository{
		collection: db.Collection("attendances"),
	}
}

func (r *attendanceRepository) Create(ctx context.Context, attendance *Attendance) error {
	attendance.ID = primitive.NewObjectID()
	attendance.CreatedAt = time.Now()
	attendance.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, &attendance)
	return err
}

func (r *attendanceRepository) FindBymEmpIDAndDate(ctx context.Context, empID primitive.ObjectID, date primitive.DateTime) (*Attendance, error) {
	var attendance Attendance
	// date, err := time.Parse("02-02-2006", dateStr)
	// if err != nil {
	// 	return nil, errors.New("invalid date format, use YYYY-MM-DD")
	// }
	// mongoDate := primitive.NewDateTimeFromTime(date)

	dberr := r.collection.FindOne(ctx, bson.M{
		"emp_id": empID,
		"date":   date,
	}).Decode(&attendance)
	if dberr != nil {
		if errors.Is(dberr, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, dberr
	}
	return &attendance, nil
}
