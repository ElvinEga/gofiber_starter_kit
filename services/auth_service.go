package services

import (
	"errors"
	"time"

	"github.com/ElvinEga/gofiber_starter/blacklist"
	"github.com/ElvinEga/gofiber_starter/config"
	"github.com/ElvinEga/gofiber_starter/database"
	"github.com/ElvinEga/gofiber_starter/models"
	"github.com/ElvinEga/gofiber_starter/requests"
	"github.com/ElvinEga/gofiber_starter/responses"
	"github.com/ElvinEga/gofiber_starter/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with auto-generated username
// @Tags Auth
// @Accept json
// @Produce json
// @Param requests.RegisterRequest body requests.RegisterRequest true "Register Request"
// @Success 201 {object} responses.AuthResponse
// @Failure 400 {object} responses.AuthResponse
// @Failure 500 {object} responses.AuthResponse
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// Check email uniqueness
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	// Create user
	newUser := models.User{
		ID:         utils.GenerateUUID(),
		Name:       req.Name,
		Email:      req.Email,
		Password:   utils.HashPassword(req.Password),
		Username:   utils.GenerateUsername(req.Name),
		Role:       "USER",
		IsVerified: false,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create user"})
	}

	token, _ := utils.GenerateJWT(newUser.ID.String())
	return c.Status(201).JSON(fiber.Map{"token": token, "user": responses.ToUserResponse(newUser)})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param requests.LoginRequest body requests.LoginRequest true "Login Request"
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.AuthResponse
// @Failure 401 {object} responses.AuthResponse
// @Failure 500 {object} responses.AuthResponse
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	var req requests.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid input",
		})
	}

	user, err := FindUserByEmail(req.Email)

	if err != nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid credentials",
		})
	}

	token, _ := utils.GenerateJWTRole(user.ID.String(), "User")
	return c.JSON(responses.AuthResponse{
		Status:  "success",
		Message: "Login successful",
		Token:   token,
		User:    responses.ToUserResponse(user),
	})
}

func generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

func GoogleLogin(c *fiber.Ctx) error {
	url := utils.GetGoogleOAuthURL()
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from Google OAuth.
func GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Code not found"})
	}

	// Exchange the code for an access token and fetch user info.
	userInfo, err := utils.GetGoogleUserInfo(code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if a user with this email exists.
	var user models.User
	if err := database.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		// If not, create a new user with auto‑generated username.
		username := utils.GenerateUsername(userInfo.Name)
		user = models.User{
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Username: username,
		}
		database.DB.Create(&user)
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})
	}

	return c.JSON(fiber.Map{"token": token, "user": user})
}

// Logout godoc
// @Summary Logout a user
// @Description Invalidate the current JWT token by blacklisting it
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} responses.AuthResponse
// @Failure 401 {object} responses.AuthResponse
// @Router /api/logout [post]
func Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Authorization header not found",
		})
	}

	// Expect token in format "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid authorization header",
		})
	}
	tokenStr := authHeader[len(bearerPrefix):]

	// Parse token to extract expiration (using the same secret).
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid token claims",
		})
	}
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(responses.AuthResponse{
			Status:  "error",
			Message: "Invalid expiration time",
		})
	}
	expirationTime := time.Unix(int64(expFloat), 0)

	// Add token to blacklist.
	blacklist.Add(tokenStr, expirationTime)

	return c.JSON(responses.AuthResponse{
		Status:  "success",
		Message: "Logout successful",
	})
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, "email = ?", email).Error
	return &user, err
}
