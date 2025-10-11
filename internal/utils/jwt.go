package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MrBista/blog-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId    int    `json:"userId"`
	Role      int    `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type JwtService struct {
	Config *config.JwtConfig
}

var jwtService *JwtService

func InitJwtService() {
	if config.AppConfig == nil {
		log.Fatal("config must be loaded before initialize JWT service")
	}
	jwtService = &JwtService{
		Config: &config.AppConfig.JWT,
	}
}

func GetJwtService() *JwtService {
	if jwtService == nil {
		log.Fatal("failed to loaded jwt services")
	}
	return jwtService
}

func (s *JwtService) CreateAccessToken(userId, role int) (string, error) {
	var secretKey = s.Config.GetSecretKey()
	expAt := s.Config.GetExpTimeAccessToken()
	claims := Claims{
		UserId:    userId,
		Role:      role,
		TokenType: "accessToken",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "blog_api",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func (s *JwtService) VerifyToken(tokenString string) (*Claims, error) {
	var secretKey = s.Config.GetSecretKey()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
