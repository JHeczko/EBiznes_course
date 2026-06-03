package middleware_handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// 🔥 FUNKCJA ZAMIAST ZMIENNEJ GLOBALNEJ: 
// Dynamicznie wyciąga klucz JWT z systemu po tym, jak .env zostanie załadowany.
func GetJwtKey() []byte {
	secret := os.Getenv("JWT_KEY")
	if secret == "" {
		// Mały bezpiecznik: jeśli zapomnisz dodać JWT_KEY do .env, aplikacja rzuci ostrzeżenie w logach
		log.Println("[WARNING] JWT_KEY is empty! Check your .env file.")
	}
	return []byte(secret)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Pobierz header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			log.Println("[AUTH ERROR] Request zablokowany: Brak nagłówka Authorization")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Brak tokena"})
		}

		// 2. Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[AUTH ERROR] Request zablokowany: Zły format nagłówka. Otrzymano: '%s' (oczekiwano 'Bearer <token>')\n", authHeader)
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Zły format tokena"})
		}
		tokenString := parts[1]

		// 3. Parsowanie i weryfikacja
		// 🔥 Zmiana: wywołujemy GetJwtKey() zamiast starej zmiennej globalnej
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return GetJwtKey(), nil
		})

		// 4. Sprawdzenie błędów (w tym wygaśnięcia!)
		if err != nil {
			log.Printf("[AUTH ERROR] Weryfikacja JWT nieudana: %v\n", err)
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token nieważny lub wygasł"})
		}

		if !token.Valid {
			log.Println("[AUTH ERROR] Request zablokowany: Token parsowany pomyślnie, ale flaga Valid = false")
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Token nieważny lub wygasł"})
		}

		// 5. Wyciągnięcie user_id z tokena i wpięcie do kontekstu
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(float64) // JWT przechowuje liczby jako float64
			if !ok {
				log.Println("[AUTH ERROR] Nie udało się wyciągnąć 'user_id' jako float64 z claims")
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Błąd struktury tokena"})
			}
			
			log.Printf("[AUTH SUCCESS] Token zweryfikowany pomyślnie. Zalogowany użytkownik ID: %d\n", int(userID))
			c.Set("user_id", int(userID))
		}

		return next(c)
	}
}