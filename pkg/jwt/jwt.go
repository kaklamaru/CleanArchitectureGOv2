package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go-clean-arch/config"
)

type JWTService struct {
	SecretKey string
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		SecretKey: cfg.JWTSecret,
	}
}

func (j *JWTService) GenerateJWT(userID uint, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)


	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err 
	}

	return signedToken, nil
}

func (j *JWTService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err 
	}

	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token") 
}
