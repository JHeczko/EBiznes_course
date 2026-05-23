package handlers

import (
	"net/http"
	"zadanie4/models"
	"zadanie4/middleware_handlers"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 3.5 Rejestracja
func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(AuthRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
		}

		// Haszowanie hasła
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
		}
		
		user := models.Users{
			Email:    req.Email,
			Password: string(hashedPassword),
			Provider: "local",
			UserName: req.Email,
		}

		if err := db.Create(&user).Error; err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "User could not be created (maybe email exists?)"})
		}

		return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
	}
}

// 3.0 Logowanie

func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(AuthRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
		}

		var user models.Users
		// 1. Znajdź usera
		if err := db.Where("email = ? AND provider = ?", req.Email, "local").First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// 2. Weryfikacja hasła
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// 3. Generowanie JWT
		expirationTime := time.Now().Add(30 * time.Minute)
		claims := jwt.MapClaims{
			"user_id": user.UserID,
			"exp":     expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(middleware_handlers.JwtKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create token"})
		}

		// 4. Odpowiedź z tokenem
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Login successful",
			"user_id": user.UserID,
			"token":   tokenString,
		})
	}
}