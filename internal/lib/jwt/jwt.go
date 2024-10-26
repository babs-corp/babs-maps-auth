package jwt

import (
	"time"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(users models.User, secret string, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = users.ID.String()
	claims["email"] = users.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
