package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the request header
		authHeader := c.Request().Header.Get("Authorization")
		// Check if the header is valid
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, "Missing auth token")
		}
		// Parse the JWT token from the header
		token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid token1")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, "Invalid token2")
		}
		// Check if the token is valid
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid token3")
		}
		// Call the next handler
		c.Set("claims", token.Claims)
		return next(c)
	}
}
