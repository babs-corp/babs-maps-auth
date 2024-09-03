package jwt

import (
	"time"

	"github.com/babs-corp/babs-maps-auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(users models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = users.ID
	claims["email"] = users.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
