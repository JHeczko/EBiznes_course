package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os" // Do wyciągania zmiennych z systemu za pomocą os.Getenv
	"time"

	"zadanie4/middleware_handlers"
	"zadanie4/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Dynamiczna konfiguracja OAuth2. Wywoływana "w locie", 
// dzięki czemu os.Getenv odpala się PO tym, jak godotenv.Load() w main.go załaduje plik .env
func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:13000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// Struktura mapująca profil, który dostajemy od API Google
type GoogleUserProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// 1. Endpoint inicjujący logowanie (React uderza tutaj, a Go przekierowuje usera do Google)
// GET /auth/google/login
func HandleGoogleLogin(c echo.Context) error {
	config := getGoogleOauthConfig()
	url := config.AuthCodeURL("random_state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// 2. Endpoint odbierający kod od Google (Callback)
// GET /auth/google/callback
func HandleGoogleCallback(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		if code == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing authorization code from Google"})
		}

		config := getGoogleOauthConfig()

		// Krok A: Wymiana jednorazowego kodu na Access Token
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to exchange token: " + err.Error()})
		}

		// Krok B: Pobranie danych profilowych użytkownika z API Google
		client := config.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user info from Google"})
		}
		defer resp.Body.Close()

		var googleUser GoogleUserProfile
		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to decode Google profile data"})
		}

		// Krok C: Sprawdzenie w bazie danych, czy użytkownik z tym mailem i providerem już istnieje
        var user models.Users
        err = db.Where("email = ?", googleUser.Email).First(&user).Error

        if err == gorm.ErrRecordNotFound {
            // Jeśli maila w ogóle nie ma w bazie, rejestrujemy nowe konto
            user = models.Users{
                Email:    googleUser.Email,
                Password: "", // Logowanie zewnętrzne, brak lokalnego hasła
                Provider: "google",
                UserName: googleUser.Name,
            }
            if err := db.Create(&user).Error; err != nil {
                return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not register user in database"})
            }
        } else if err == nil { // jeśli jest w bazie, wtedy nalezy zmienic providera, na najnowszego
            if user.Provider != "google" {
                if err := db.Model(&user).Update("provider", "google").Error; err != nil {
                    return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not update user provider"})
                }
            }
        } else {
            return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database query error"})
        }

		// Krok D: Generowanie Twojego wewnętrznego JWT (dla Twojego AuthMiddleware)
		expirationTime := time.Now().Add(30 * time.Minute)
		claims := jwt.MapClaims{
			"user_id": user.UserID,
			"exp":     expirationTime.Unix(),
		}

		appToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		appTokenString, err := appToken.SignedString(middleware_handlers.GetJwtKey())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate application token"})
		}

		// Krok E: Przekierowanie z powrotem do Reacta z tokenem i userID w parametrach URL
		reactRedirectURL := fmt.Sprintf("http://localhost:5173/login?token=%s&user_id=%d&provider=google", appTokenString, user.UserID)
		return c.Redirect(http.StatusMovedPermanently, reactRedirectURL)
	}
}