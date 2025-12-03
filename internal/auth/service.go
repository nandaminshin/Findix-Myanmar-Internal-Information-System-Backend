package auth

import (
	"errors"
	"time"

	"fmiis/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type authService struct {
	jwtSecret []byte
}

func NewAuthService() AuthService {
	cfg, _ := config.LoadConfig()
	return &authService{
		jwtSecret: []byte(cfg.JWTSecret),
	}
}

func (s *authService) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Second * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
}
