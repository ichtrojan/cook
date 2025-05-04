package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/ichtrojan/cook/config"
	"github.com/ichtrojan/cook/database"
	"github.com/ichtrojan/cook/helpers"
	"github.com/ichtrojan/cook/models"
	"net/http"
	"strings"
)

type userCtxKey string

const (
	userKey userCtxKey = "user"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)

		unauthorized := helpers.PrepareMessage("unauthorized")

		if token == "" {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusUnauthorized)

			_ = json.NewEncoder(w).Encode(unauthorized)

			return
		}

		validation, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.AppConfig.AppKey), nil
		})

		var foundUser models.User

		if claims, ok := validation.Claims.(jwt.MapClaims); ok && validation.Valid {
			tempToken, exist := claims["token"].(string)

			if !exist {
				tempToken = ""
			}

			userId := database.Redis.Get(r.Context(), "user_auth_"+tempToken).Val()

			_ = database.Mysql.Where("id = ?", userId).First(&foundUser)

			if foundUser.Empty() {
				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(http.StatusUnauthorized)

				_ = json.NewEncoder(w).Encode(unauthorized)

				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userKey, foundUser))

			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusUnauthorized)

			_ = json.NewEncoder(w).Encode(unauthorized)

			return
		}
	})
}

func GetUser(ctx context.Context) models.User {
	return ctx.Value(userKey).(models.User)
}

func IsUser(ctx context.Context, user models.User) bool {
	return ctx.Value(userKey).(models.User).Id == user.Id
}
