package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/bobykurniawan11/starter-go-prisma/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(id string) (string, error) {

	config := config.GetConfig()
	secret := config.GetString("api.secret")
	token_lifespan := config.GetInt("token.life")

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func TokenValid(c *gin.Context) error {

	config := config.GetConfig()
	secret := config.GetString("api.secret")

	tokenString := ExtractToken(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	return nil
}

// func ExtractTokenID(c *gin.Context) (uuid.UUID, error) {

// 	config := config.GetConfig()
// 	secret := config.GetString("api.secret")
// 	tokenString := ExtractToken(c)
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secret), nil

// 	})
// 	if err != nil {
// 		return uuid.UUID{}, err
// 	}
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if ok && token.Valid {
// 		uid := claims["user_id"]

// 		return uid.(uuid.UUID), nil
// 	}
// 	return uuid.UUID{}, nil

// }

func ExtractTokenID(c *gin.Context) (uuid.UUID, error) {
	config := config.GetConfig()
	secret := config.GetString("api.secret")
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid := claims["user_id"]

		// Convert uid to string
		uidStr, ok := uid.(string)
		if !ok {
			return uuid.UUID{}, fmt.Errorf("Failed to convert user_id to string")
		}

		// Parse string to uuid.UUID
		uidUUID, err := uuid.Parse(uidStr)
		if err != nil {
			return uuid.UUID{}, err
		}

		return uidUUID, nil
	}
	return uuid.UUID{}, nil
}
