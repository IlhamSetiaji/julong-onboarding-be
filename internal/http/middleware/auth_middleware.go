package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-onboarding-be/internal/http/request"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewAuth(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("auth", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUser(c *gin.Context, log *logrus.Logger) (map[string]interface{}, error) {
	auth, exists := c.Get("auth")
	if !exists {
		return nil, errors.New("auth key not found in context")
	}

	claims, ok := auth.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid auth claims type")
	}

	message := messaging.UserMessageFactory(log)

	messageResponse, err := message.SendGetUserMe(request.SendFindUserByIDMessageRequest{
		ID: claims["id"].(string),
	})

	if err != nil {
		log.Errorf("[GetUserMe.Create Message] " + err.Error())
		return nil, err
	}

	if messageResponse.User == nil {
		log.Errorf("[GetUserMe.Create] User not found")
		return nil, errors.New("User not found")
	}

	return messageResponse.User, nil
}

func CheckLoggedIn(c *gin.Context, log *logrus.Logger) bool {
	_, exists := c.Get("auth")
	if !exists {
		return false
	}

	return true
}
