package user

import (
	"context"
	"errors"
	"fmiis/internal/auth"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GmUpdate(ctx context.Context, req *GmUpdateRequest) (*UserResponse, error)
	NormalUpdate(ctx context.Context, req *NormalUpdateRequest) (*UserResponse, error)
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Default role to "user" if not provided
	role := req.Role
	if role == "" {
		role = "glob"
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
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
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Image != "" {
		user.Image = req.Image
	}

	// Optional: update password
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.Update(ctx, user.ID, user); err != nil {
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
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Image != "" {
		user.Image = req.Image
	}

	if err := s.repo.Update(ctx, user.ID, user); err != nil {
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
