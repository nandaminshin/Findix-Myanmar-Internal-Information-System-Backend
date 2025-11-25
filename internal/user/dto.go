package user

type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=3"`
	Phone      string `json:"phone" binding:"required"`
	Role       Role   `json:"role" binding:"required"`
	SecretCode string `json:"secret_code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  Role   `json:"role"`
	Image string `json:"image,omitempty"`
	Token string `json:"token,omitempty"`
}
