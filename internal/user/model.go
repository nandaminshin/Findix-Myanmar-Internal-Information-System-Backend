package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	dev  Role = "dev"
	glob Role = "glob"
	hr   Role = "hr"
	gm   Role = "gm"
	md   Role = "md"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	Password         string             `bson:"password" json:"-"`
	Role             Role               `bson:"role" json:"role"`
	Image            string             `bson:"image,omitempty" json:"image,omitempty"`
	EmpNumber        string             `bson:"emp_no" json:"emp_no"`
	Birthday         primitive.DateTime `bson:"birthday" json:"birthday"`
	DateOfHire       primitive.DateTime `bson:"date_of_hire" json:"date_of_hire"`
	Salary           int64              `bson:"salary" json:"salary"`
	DateOfRetirement primitive.DateTime `bson:"date_of_retirement" json:"date_of_retirement"`
	NRC              string             `bson:"nrc" json:"nrc"`
	GraduatedUni     string             `bson:"graduated_uni" json:"graduated_uni"`
	Address          string             `bson:"address" json:"address"`
	ParentAddress    string             `bson:"parent_address" json:"parent_address"`
	ParentPhone      string             `bson:"parent_phone" json:"parent_phone"`
	Note             string             `bson:"note" json:"note"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
