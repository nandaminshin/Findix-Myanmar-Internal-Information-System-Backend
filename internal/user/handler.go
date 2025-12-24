package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    res.ID,
		"name":  res.Name,
		"email": res.Email,
		"role":  res.Role,
		"phone": res.Phone,
		"image": res.Image,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt_token", res.Token, 168*3600, "/", "", true, true)

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *UserHandler) GmUpdate(c *gin.Context) {
	var req GmUpdateRequest
	userID := c.Param("id")
	user, err := h.service.GetSingleUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authUserRole := c.MustGet("role").(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authUserRole == string(Hr) && (user.Role == Gm || user.Role == Md) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this user"})
		return
	}
	res, err := h.service.GmUpdate(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":  res.Name,
		"email": res.Email,
		"role":  res.Role,
		"phone": res.Phone,
		"image": res.Image,
	})
}

func (h *UserHandler) NormalUpdate(c *gin.Context) {
	var req NormalUpdateRequest
	userID := c.Param("id")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.service.NormalUpdate(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":  res.Name,
		"image": res.Image,
	})
}

func (h *UserHandler) GmDelete(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.DeleteById(c, userID, req.SecretCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *UserHandler) GetAllEmployees(c *gin.Context) {
	res, err := h.service.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"employees": res,
	})
}

func (h *UserHandler) GetSingleEmployee(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	res, err := h.service.GetSingleUser(c, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"employee": res,
	})
}

func (h *UserHandler) UploadProfileImage(c *gin.Context) {
	file, err := c.FormFile("profile_image")
	if err != nil {
		c.JSON(400, gin.H{"error": "profile_picture is required"})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// Call the unified service method
	res, err := h.service.UpdateProfileImage(c.Request.Context(), userID, file)
	if err != nil {
		code := http.StatusInternalServerError
		if err.Error() == "user not found" || err.Error() == "file size too large" || strings.Contains(err.Error(), "file type not allowed") || strings.Contains(err.Error(), "invalid image type") {
			code = http.StatusBadRequest
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"updated_employee": res,
		"message":          "upload success",
		"path":             res.Image,
	})
}
