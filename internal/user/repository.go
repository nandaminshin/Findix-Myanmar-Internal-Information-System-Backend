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
	FetchAllUsers(ctx context.Context) (*[]User, error)
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
	case Dev, Glob, Hr, Gm, Md:
		return true
	default:
		return false
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if !IsValidRole(user.Role) {
		return errors.New("invalid role")
	}

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

func (r *mongoUserRepository) FetchAllUsers(ctx context.Context) (*[]User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, id primitive.ObjectID, u *User) error {
	if !IsValidRole(u.Role) {
		return errors.New("invalid role")
	}

	if u.Image != "" {
		_, err := r.collection.UpdateByID(ctx, id,
			bson.M{
				"$set": bson.M{
					"name":               u.Name,
					"email":              u.Email,
					"image":              u.Image,
					"phone":              u.Phone,
					"role":               u.Role,
					"emp_no":             u.EmpNumber,
					"birthday":           u.Birthday,
					"date_of_hire":       u.DateOfHire,
					"salary":             u.Salary,
					"date_of_retirement": u.DateOfRetirement,
					"nrc":                u.NRC,
					"graduated_uni":      u.GraduatedUni,
					"address":            u.Address,
					"emergency_address":  u.EmergencyAddress,
					"emergency_phone":    u.EmergencyPhone,
					"family_info":        u.FamilyInfo,
					"leave_info":         u.LeaveInfo,
					"note":               u.Note,
					"password":           u.Password,
				},
			},
		)
		return err
	} else {
		_, err := r.collection.UpdateByID(ctx, id,
			bson.M{
				"$set": bson.M{
					"name":               u.Name,
					"email":              u.Email,
					"phone":              u.Phone,
					"role":               u.Role,
					"emp_no":             u.EmpNumber,
					"birthday":           u.Birthday,
					"date_of_hire":       u.DateOfHire,
					"salary":             u.Salary,
					"date_of_retirement": u.DateOfRetirement,
					"nrc":                u.NRC,
					"graduated_uni":      u.GraduatedUni,
					"address":            u.Address,
					"emergency_address":  u.EmergencyAddress,
					"emergency_phone":    u.EmergencyPhone,
					"family_info":        u.FamilyInfo,
					"leave_info":         u.LeaveInfo,
					"note":               u.Note,
					"password":           u.Password,
				},
			},
		)
		return err
	}
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
