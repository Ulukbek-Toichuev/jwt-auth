package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User interface {
	GetUserName() string
	GetEmail() string
	GetRole() string
	GetPassword() string
}

func GenerateToken(user User, secretKey string, tokenExpiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userName": user.GetUserName(),
		"email":    user.GetEmail(),
		"roles":    user.GetRole(),
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(tokenExpiry).Unix(),
	})

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", jwt.ErrTokenSignatureInvalid
	}

	return signedToken, nil
}

func VerifyToken(token, secretKey string) (map[string]interface{}, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
