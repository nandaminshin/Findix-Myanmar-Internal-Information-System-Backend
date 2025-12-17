package user

import (
	"context"
	"errors"
	"fmiis/internal/auth"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GmUpdate(ctx context.Context, req *GmUpdateRequest) (*UserResponse, error)
	NormalUpdate(ctx context.Context, req *NormalUpdateRequest) (*UserResponse, error)
	GetAllUsers(ctx context.Context) (*[]User, error)
	GetSingleUser(ctx context.Context, id string) (*User, error)
	DeleteById(ctx context.Context, userID string, secretCode string) error
	UpdateProfileImage(ctx context.Context, userID string, file *multipart.FileHeader) (*UserResponse, error)
}

type userService struct {
	repo          UserRepository
	authService   auth.AuthService
	socketServer  socketio.Server
	uploadBaseDir string
}

func NewUserService(repo UserRepository, authService auth.AuthService, server *socketio.Server, uploadBaseDir string) UserService {
	return &userService{
		repo:          repo,
		authService:   authService,
		socketServer:  *server,
		uploadBaseDir: uploadBaseDir,
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
		Name:             req.Name,
		Email:            req.Email,
		Password:         string(hashedPassword),
		Phone:            req.Phone,
		Role:             role,
		EmpNumber:        req.EmpNumber,
		Birthday:         birthday,
		DateOfHire:       dateOfHire,
		Salary:           req.Salary,
		NRC:              req.NRC,
		GraduatedUni:     req.GraduatedUni,
		Address:          req.Address,
		EmergencyAddress: req.EmergencyAddress,
		EmergencyPhone:   req.EmergencyPhone,
		FamilyInfo:       req.FamilyInfo,
		Note:             req.Note,
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

	// ... database creation ...

	// WRAP THIS IN A GOROUTINE
	go func() {
		s.socketServer.BroadcastToNamespace(
			"/",
			"employee_created",
			gin.H{
				"message": "New employee created",
			},
		)
	}()

	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
		Image: user.Image,
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

func (s *userService) GetAllUsers(ctx context.Context) (*[]User, error) {
	fetchedUsers, err := s.repo.FetchAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	return fetchedUsers, err
}

func (s *userService) GetSingleUser(ctx context.Context, userID string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	user, err := s.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GmUpdate(ctx context.Context, req *GmUpdateRequest) (*UserResponse, error) {
	if req.SecretCode != os.Getenv("SECRET_CODE") {
		return nil, errors.New("access denied, invalid secret code")
	}

	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.FindByID(ctx, objID)
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
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		Role:             req.Role,
		EmpNumber:        req.EmpNumber,
		Birthday:         birthday,
		DateOfHire:       dateOfHire,
		Salary:           req.Salary,
		NRC:              req.NRC,
		GraduatedUni:     req.GraduatedUni,
		Address:          req.Address,
		EmergencyAddress: req.EmergencyAddress,
		EmergencyPhone:   req.EmergencyPhone,
		FamilyInfo:       req.FamilyInfo,
		Note:             req.Note,
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

func (s *userService) DeleteById(ctx context.Context, userID string, secretCode string) error {
	if secretCode != os.Getenv("SECRET_CODE") {
		return errors.New("access denied, invalid secret code")
	}

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.repo.Delete(ctx, objID)
}

func (s *userService) ParseDateTimeForDB(dt string) (primitive.DateTime, error) {
	t, err := time.Parse("2006-01-02", dt)
	if err != nil {
		return 0, errors.New("invalid date format, expected dd-mm-yy")
	}
	// Convert to primitive.DateTime for MongoDB
	dbDateTime := primitive.NewDateTimeFromTime(t)
	return dbDateTime, nil
}

func (s *userService) validateImage(file *multipart.FileHeader) error {
	const maxSize = 10 * 1024 * 1024
	if file.Size > maxSize {
		return errors.New("file size too large")
	}
	// 2. extension check
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	if !allowedExt[ext] {
		return errors.New("file type not allowed, only jpg, jpeg and png are allowed")
	}

	// 3. MIME type check
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	if _, err := src.Read(buffer); err != nil {
		return err
	}

	// You might want to add these checks:
	if len(buffer) == 0 {
		return errors.New("empty file")
	}

	contentType := http.DetectContentType(buffer)
	allowedMime := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !allowedMime[contentType] {
		return errors.New("invalid image type")
	}

	// Reset file pointer if you need to process the file further
	if _, err := src.Seek(0, 0); err != nil {
		return err
	}
	return nil
}

func (s *userService) saveProfilePicture(file *multipart.FileHeader, employeeID string) (string, error) {
	uploadDir := filepath.Join(s.uploadBaseDir, "employee", "images")

	// ensure folder exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	// safer filename
	filename := fmt.Sprintf("fmiis_user_%s%s", employeeID, ext)
	fullFilePath := filepath.Join(uploadDir, filename)

	if err := SaveMultipartFile(file, fullFilePath); err != nil {
		return "", err
	}

	return fullFilePath, nil
}

func SaveMultipartFile(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (s *userService) UpdateProfileImage(ctx context.Context, userID string, file *multipart.FileHeader) (*UserResponse, error) {
	// 1. Validate file
	if err := s.validateImage(file); err != nil {
		return nil, err
	}

	// 2. Get user
	user, err := s.GetSingleUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 3. Delete old image if exists
	if user.Image != "" {
		// Resolve path for deletion
		var pathToDelete string
		// Check if it's our web-relative path
		if strings.HasPrefix(user.Image, "/uploads/") {
			// Strip /uploads prefix and join with actual base dir
			cleanPath := strings.TrimPrefix(user.Image, "/uploads")
			pathToDelete = filepath.Join(s.uploadBaseDir, cleanPath)
		} else {
			// Assume it's a legacy absolute path or other format
			pathToDelete = user.Image
		}

		if err := os.Remove(pathToDelete); err != nil {
			fmt.Printf("Failed to delete old image %s: %v\n", pathToDelete, err)
		}
	}

	// 4. Save new file
	fullFilePath, err := s.saveProfilePicture(file, userID)
	if err != nil {
		return nil, err
	}

	// Construct relative path for DB/Frontend
	// fullFilePath is like /path/to/app/uploads/employee/images/file.png
	// We want /uploads/employee/images/file.png
	// s.uploadBaseDir is /path/to/app/uploads

	// Ensure uniform separators
	relPath, err := filepath.Rel(s.uploadBaseDir, fullFilePath)
	if err != nil {
		// Fallback to full path or error?
		// If Rel fails, something weird. Just use full path as fallback or log.
		// Usually won't fail if fullFilePath is inside uploadBaseDir.
		relPath = filepath.Base(fullFilePath) // Too risky?
	}

	// filepath.Rel returns "employee/images/file.png" (no leading slash)
	// We want "/uploads/" + relPath
	// Use slashes for URL compatibility
	urlPath := "/uploads/" + filepath.ToSlash(relPath)

	// 5. Update user in DB
	user.Image = urlPath
	if err := s.repo.Update(ctx, user.ID, user); err != nil {
		return nil, err
	}

	// 6. Return response
	return &UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
		Phone: user.Phone,
		Image: user.Image,
	}, nil
}
