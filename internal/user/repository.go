package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByEmpNo(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	Update(ctx context.Context, id primitive.ObjectID, user *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func IsValidRole(role Role) bool {
	switch role {
	case "dev", "glob", "hr", "gm", "md":
		return true
	default:
		return false
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, &user)
	return err
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) FindByEmpNo(ctx context.Context, empNo string) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"emp_no": empNo}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, id primitive.ObjectID, u *User) error {
	_, err := r.collection.UpdateByID(ctx, id,
		bson.M{
			"$set": bson.M{
				"name":             u.Name,
				"email":            u.Email,
				"phone":            u.Phone,
				"role":             u.Role,
				"empNumber":        u.EmpNumber,
				"birthday":         u.Birthday,
				"dateOfHire":       u.DateOfHire,
				"salary":           u.Salary,
				"dateOfRetirement": u.DateOfRetirement,
				"nrc":              u.NRC,
				"graduatedUni":     u.GraduatedUni,
				"address":          u.Address,
				"parentAddress":    u.ParentAddress,
				"parentPhone":      u.ParentPhone,
				"note":             u.Note,
				"password":         u.Password,
			},
		},
	)
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
