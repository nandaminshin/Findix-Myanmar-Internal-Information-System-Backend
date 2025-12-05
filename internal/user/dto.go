package user

type RegisterRequest struct {
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=3"`
	Phone            string `json:"phone" binding:"required"`
	Role             Role   `json:"role" binding:"required"`
	EmpNumber        string `json:"emp_no" binding:"required"`
	Birthday         string `json:"birthday" binding:"required"`
	DateOfHire       string `json:"date_of_hire" binding:"required"`
	Salary           int64  `json:"salary" binding:"required"`
	DateOfRetirement string `json:"date_of_retirement"`
	NRC              string `json:"nrc" binding:"required"`
	GraduatedUni     string `json:"graduated_uni"`
	Address          string `json:"address" binding:"required"`
	ParentAddress    string `json:"parent_address" binding:"required"`
	ParentPhone      string `json:"parent_phone" binding:"required"`
	Note             string `json:"note"`
	SecretCode       string `json:"secret_code" binding:"required"`
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

type GmUpdateRequest struct {
	ID               string `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=3"`
	Phone            string `json:"phone" binding:"required"`
	Role             Role   `json:"role" binding:"required"`
	EmpNumber        string `json:"emp_no" binding:"required"`
	Birthday         string `json:"birthday" binding:"required"`
	DateOfHire       string `json:"date_of_hire" binding:"required"`
	Salary           int64  `json:"salary" binding:"required"`
	DateOfRetirement string `json:"date_of_retirement"`
	NRC              string `json:"nrc" binding:"required"`
	GraduatedUni     string `json:"graduated_uni"`
	Address          string `json:"address" binding:"required"`
	ParentAddress    string `json:"parent_address" binding:"required"`
	ParentPhone      string `json:"parent_phone" binding:"required"`
	Note             string `json:"note"`
	SecretCode       string `json:"secret_code" binding:"required"`
	Image            string `json:"image,omitempty"`
}

type NormalUpdateRequest struct {
	ID    string `json:"id" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Image string `json:"image,omitempty"`
}

type DeleteRequest struct {
	ID string `json:"id" binding:"required"`
}
