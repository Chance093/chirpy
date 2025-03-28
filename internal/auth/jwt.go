package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, duration time.Duration) (string, error) {
	currentTime := time.Now().UTC()
	expirationTime := currentTime.Add(duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   userID.String(),
		},
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error while signing json token: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uuid.Nil, fmt.Errorf("error while parsing claims from token: %v", err)
	}

  if claims.Issuer != "chirpy" {
    return uuid.Nil, errors.New("invalid issuer")
  }

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while parsing userId string into uuid: %v", err)
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
  authHeader := headers.Get("Authorization")
  if authHeader == "" {
    return "", errors.New("Authorization header missing")
  }

  if !strings.Contains(authHeader, "Bearer ") {
    return "", errors.New("Authorization header missing Bearer token")
  }

  token := strings.TrimPrefix(authHeader, "Bearer ")

  return token, nil 
}
