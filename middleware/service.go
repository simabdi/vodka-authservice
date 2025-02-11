package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/simabdi/vodka-authservice/config"
	"github.com/simabdi/vodka-authservice/helper"
	"github.com/simabdi/vodka-authservice/models"
	"strconv"
	"time"
)

type Service interface {
	GenerateToken(user models.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	VerifyToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewJwtService() *jwtService {
	return &jwtService{}
}

func (js *jwtService) GenerateToken(user models.User) (string, error) {
	lifeTime, _ := strconv.Atoi(config.LifeTime)
	ttl := time.Duration(lifeTime) * time.Second
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid":  user.Uuid,
		"refId": user.RefID,
		"role":  user.Role,
		"exp":   time.Now().UTC().Add(ttl).Unix(),
	})

	signedToken, err := token.SignedString([]byte(helper.Std64Decode(config.JWTSecretKey)))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (js *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	resultToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte("secret_key"), nil
	})

	if err != nil {
		return resultToken, err
	}

	return resultToken, nil
}

func (js *jwtService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(helper.Std64Decode(config.JWTSecretKey)), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
