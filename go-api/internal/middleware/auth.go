package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/xangel1111/qr-stats-fullstack/go-api/internal/apperr"
)

// Locals keys used to pass request-scoped values to handlers.
const (
	UsernameKey = "username"
	TokenKey    = "token"
)

// GenerateToken issues an HS256 JWT for the given subject (username).
func GenerateToken(secret, username string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": username,
		"iat": now.Unix(),
		"exp": now.Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTAuth validates the Bearer token and stores the subject and raw token in
// the request locals. The raw token is propagated to the Node service.
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return apperr.Unauthorized("missing or malformed Authorization header")
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return apperr.Unauthorized("invalid or expired token")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sub, ok := claims["sub"].(string); ok {
				c.Locals(UsernameKey, sub)
			}
		}
		c.Locals(TokenKey, tokenStr)
		return c.Next()
	}
}
