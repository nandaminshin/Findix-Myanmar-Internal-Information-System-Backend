package user

import (
	"context"
	"errors"
	"fmiis/internal/auth"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GmUpdate(ctx context.Context, req *GmUpdateRequest) (*UserResponse, error)
	NormalUpdate(ctx context.Context, req *NormalUpdateRequest) (*UserResponse, error)
	DeleteById(ctx context.Context, userID string) error
}

type userService struct {
	repo        UserRepository
	authService auth.AuthService
}

func NewUserService(repo UserRepository, authService auth.AuthService) UserService {
	return &userService{
		repo:        repo,
		authService: authService,
	}
}

func (s *userService) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	if req.SecretCode != os.Getenv("SECRET_CODE") {
		return nil, errors.New("access denied, invalid secret code")
	}
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUserByEmpNo, err := s.repo.FindByEmpNo(ctx, req.EmpNumber)
	if err != nil {
		return nil, err
	}
	if existingUserByEmpNo != nil {
		return nil, errors.New("employee number already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Default role to "user" if not provided
	role := req.Role
	if role == "" {
		role = "glob"
	}

	birthday, err := s.ParseDateTimeForDB(req.Birthday)
	if err != nil {
		return nil, err
	}

	dateOfHire, err := s.ParseDateTimeForDB(req.DateOfHire)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:          req.Name,
		Email:         req.Email,
		Password:      string(hashedPassword),
		Phone:         req.Phone,
		Role:          role,
		EmpNumber:     req.EmpNumber,
		Birthday:      birthday,
		DateOfHire:    dateOfHire,
		Salary:        req.Salary,
		NRC:           req.NRC,
		GraduatedUni:  req.GraduatedUni,
		Address:       req.Address,
		ParentAddress: req.ParentAddress,
		ParentPhone:   req.ParentPhone,
		Note:          req.Note,
	}

	if req.DateOfRetirement != "" {
		dateOfRetirement, err := s.ParseDateTimeForDB(req.DateOfRetirement)
		if err != nil {
			return nil, err
		}
		user.DateOfRetirement = dateOfRetirement
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
	}, nil
}

func (s *userService) Login(ctx context.Context, req *LoginRequest) (*UserResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.authService.GenerateToken(user.ID.Hex(), string(user.Role))
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Phone: user.Phone,
		Image: user.Image,
		Token: token,
	}, nil

}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("Error getting user by email %s: %v", email, err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GmUpdate(ctx context.Context, req *GmUpdateRequest) (*UserResponse, error) {
	if req.SecretCode != os.Getenv("SECRET_CODE") {
		return nil, errors.New("access denied, invalid secret code")
	}

	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Email != user.Email {
		otherUser, err := s.repo.FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if otherUser != nil && otherUser.ID != user.ID {
			return nil, errors.New("email already exists")
		}
	}

	birthday, err := s.ParseDateTimeForDB(req.Birthday)
	if err != nil {
		return nil, err
	}

	dateOfHire, err := s.ParseDateTimeForDB(req.DateOfHire)
	if err != nil {
		return nil, err
	}

	updatedUser := &User{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Role:          req.Role,
		EmpNumber:     req.EmpNumber,
		Birthday:      birthday,
		DateOfHire:    dateOfHire,
		Salary:        req.Salary,
		NRC:           req.NRC,
		GraduatedUni:  req.GraduatedUni,
		Address:       req.Address,
		ParentAddress: req.ParentAddress,
		ParentPhone:   req.ParentPhone,
		Note:          req.Note,
	}

	if req.DateOfRetirement != "" {
		dateOfRetirement, err := s.ParseDateTimeForDB(req.DateOfRetirement)
		if err != nil {
			return nil, err
		}
		updatedUser.DateOfRetirement = dateOfRetirement
	}

	// Optional: update password
	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if passwordErr == nil {
		// password matches existing → not changed
		updatedUser.Password = user.Password
	} else {
		// password changed → hash new one
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updatedUser.Password = string(hashed)
	}

	if err := s.repo.Update(ctx, user.ID, updatedUser); err != nil {
		return nil, err
	}

	// 5. Return updated response
	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Phone: user.Phone,
		Image: user.Image,
	}, nil
}

func (s *userService) NormalUpdate(ctx context.Context, req *NormalUpdateRequest) (*UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	updatedUser := &User{
		Name:  req.Name,
		Image: req.Image,
	}

	if err := s.repo.Update(ctx, user.ID, updatedUser); err != nil {
		return nil, err
	}

	// 5. Return updated response
	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Phone: user.Phone,
		Image: user.Image,
	}, nil
}

func (s *userService) DeleteById(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(ctx, objID)
}

func (s *userService) ParseDateTimeForDB(dt string) (primitive.DateTime, error) {
	t, err := time.Parse("02/01/2006", dt) //dd/mm/yyyy format date
	if err != nil {
		return 0, errors.New("invalid date format, expected dd/mm/yy")
	}
	// Convert to primitive.DateTime for MongoDB
	dbDateTime := primitive.NewDateTimeFromTime(t)
	return dbDateTime, nil
}
