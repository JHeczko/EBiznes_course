package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"zadanie4/middleware_handlers"
	"zadanie4/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

func getGithubOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:13000/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLEINT_SECRET"),
		Scopes:       []string{"user:email", "read:user"},
		Endpoint:     github.Endpoint,
	}
}

// Struktura dla danych profilowych z GitHuba
type GithubUserProfile struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Struktura pomocnicza, bo GitHub maile zwraca w osobnym endpoincie jako tablicę
type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// 1. Inicjalizacja logowania przez GitHub
// GET /auth/github/login
func HandleGithubLogin(c echo.Context) error {
	config := getGithubOauthConfig()
	url := config.AuthCodeURL("random_github_state")
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// 2. Callback z GitHuba
// GET /auth/github/callback
func HandleGithubCallback(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		if code == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Missing authorization code from GitHub"})
		}

		config := getGithubOauthConfig()

		// Krok A: Wymiana kodu na Access Token
		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to exchange GitHub token: " + err.Error()})
		}

		client := config.Client(context.Background(), token)

		// Krok B1: Pobranie podstawowego profilu
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user profile from GitHub"})
		}
		defer resp.Body.Close()

		var githubUser GithubUserProfile
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to decode GitHub profile"})
		}

		// Krok B2: Ponieważ profil w GitHubie może mieć ukryty mail, uderzamy po listę maili użytkownika
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []GithubEmail
			if err := json.NewDecoder(emailResp.Body).Decode(&emails); err == nil {
				for _, e := range emails {
					if e.Primary { // Szukamy głównego maila
						githubUser.Email = e.Email
						break
					}
				}
			}
		}

		// Zabezpieczenie awaryjne, gdyby user nie miał publicznego maila
		if githubUser.Email == "" {
			githubUser.Email = fmt.Sprintf("%s@github-user.com", githubUser.Login)
		}
		if githubUser.Name == "" {
			githubUser.Name = githubUser.Login
		}

		// Krok C: Zapisanie danych logowania OAuth2 po stronie serwera 
		var user models.Users
		err = db.Where("email = ?", githubUser.Email).First(&user).Error

		if err == gorm.ErrRecordNotFound {
			// Jeśli maila w ogóle nie ma w bazie, rejestrujemy nowego usera
			user = models.Users{
				Email:    githubUser.Email,
				Password: "",
				Provider: "github",
				UserName: githubUser.Name,
			}
			if err := db.Create(&user).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not register user"})
			}
		} else if err == nil {
			// 🔥 Jeśli user istnieje, ale zalogował się innym providerem, 
			// możemy opcjonalnie zaktualizować providera na najnowszy lub po prostu wpuścić usera:
			if user.Provider != "github" {
				db.Model(&user).Update("provider", "github")
			}
		} else {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		}

		// Krok D: Generowanie Twojego wewnętrznego JWT (z użyciem poprawionej funkcji GetJwtKey!)
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

		// Krok E: Przekierowanie do Reacta
		reactRedirectURL := fmt.Sprintf("http://localhost:5173/login?token=%s&user_id=%d&provider=github", appTokenString, user.UserID)
		return c.Redirect(http.StatusMovedPermanently, reactRedirectURL)
	}
}