package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FiodhyAN/learn-rest-api-golang/config"
	"github.com/FiodhyAN/learn-rest-api-golang/helpers"
	"github.com/FiodhyAN/learn-rest-api-golang/types"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func CreateJWT(secret []byte, userID string) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestToken := getTokenFromRequest(r)

		token, err := validateToken(requestToken)
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID := str

		u, err := store.GetUserById(userID)
		if err != nil {
			log.Printf("failed to get user id: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token != "" {
		return token
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	helpers.WriteError(w, http.StatusForbidden, "Forbidden", fmt.Errorf("Permission Denied"))
}

func GetUserIdFromContext(ctx context.Context) string {
	userID, ok := ctx.Value(UserKey).(string)
	if !ok {
		return ""
	}

	return userID
}
