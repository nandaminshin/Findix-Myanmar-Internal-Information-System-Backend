package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	Dev  Role = "dev"
	Glob Role = "glob"
	Hr   Role = "hr"
	Gm   Role = "gm"
	Md   Role = "md"
)

type FamilyInfo struct {
	Dad             bool  `bson:"dad" json:"dad"`
	DadAllowance    bool  `bson:"dad_allowance" json:"dad_allowance"`
	Mom             bool  `bson:"mom" json:"mom"`
	MomAllowance    bool  `bson:"mom_allowance" json:"mom_allowance"`
	SpouseAllowance bool  `bson:"spouse_allowance" json:"spouse_allowance"`
	Child           int16 `bson:"child" json:"child"`
}

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
	EmergencyAddress string             `bson:"emergency_address" json:"emergency_address"`
	EmergencyPhone   string             `bson:"emergency_phone" json:"emergency_phone"`
	FamilyInfo       FamilyInfo         `bson:"family_info" json:"family_info"`
	Note             string             `bson:"note" json:"note"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
