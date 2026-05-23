package handlers

import (
	"net/http"
	"zadanie4/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		// Szukamy użytkownika lokalnego po emailu
		if err := db.Where("Email = ? AND Provider = ?", req.Email, "local").First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// Weryfikacja hasła
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
		}

		// Tu w przyszłości wstawisz JWT - teraz tylko info, że zalogowany
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Login successful", 
			"user_id": user.UserID,
		})
	}
}