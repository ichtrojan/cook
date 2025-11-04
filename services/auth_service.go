package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/ichtrojan/cook/config"
	"github.com/ichtrojan/cook/database"
	"github.com/ichtrojan/cook/helpers"
	"github.com/ichtrojan/cook/mailer"
	middleware "github.com/ichtrojan/cook/middlewares"
	"github.com/ichtrojan/cook/models"
	"github.com/ichtrojan/cook/queue"
	"github.com/ichtrojan/cook/requests"
	"github.com/ichtrojan/cook/responses"
	"github.com/ichtrojan/gotp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(payload requests.Register) (response responses.AuthResponse, err error, status int) {
	var user models.User

	_ = database.Mysql.Where("email = ?", payload.Email).Find(&user)

	if !user.Empty() {
		return response, errors.New("email already taken"), http.StatusNotAcceptable
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return responses.AuthResponse{}, helpers.ServerError(err), http.StatusInternalServerError
	}

	user = models.User{
		Id:        uuid.New().String(),
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = database.Mysql.Create(&user).Error; err != nil {
		return response, helpers.ServerError(err), http.StatusInternalServerError
	}

	tokenStore := uuid.New().String()

	tokenExpiry := time.Hour * (24 * 7)

	err = database.Redis.Set(
		context.Background(),
		"user_auth_"+tokenStore,
		user.Id,
		tokenExpiry,
	).Err()

	if err != nil {
		return responses.AuthResponse{}, helpers.ServerError(err), http.StatusInternalServerError
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(tokenExpiry).Unix(),
		"token": tokenStore,
	})

	token, err := at.SignedString([]byte(config.AppConfig.AppKey))

	if err != nil {
		return response, helpers.ServerError(err), http.StatusInternalServerError
	}

	otp, err := gotp.New(gotp.Config{Redis: database.Redis})

	if err != nil {
		return response, helpers.ServerError(err), http.StatusInternalServerError
	}

	otpToken, err := otp.Generate(gotp.Generate{
		Format:     gotp.NUMERIC,
		Length:     6,
		Identifier: fmt.Sprintf("signup_%s", user.Id),
		Expires:    time.Minute * 10,
	})

	if err != nil {
		return response, helpers.ServerError(err), http.StatusInternalServerError
	}

	err = mailer.EnqueueEmailTask(queue.Client, mailer.EmailPayload{
		TemplateName: "signup_otp",
		To:           user.Email,
		Subject:      "OTP from Cook",
		Data: map[string]interface{}{
			"Token": otpToken,
			"Name":  user.Name,
			"Email": user.Email,
		},
	})

	if err != nil {
		return response, helpers.ServerError(err), http.StatusInternalServerError
	}

	return responses.AuthResponse{
		Token: token,
		User:  responses.GenerateUserResponse(user),
	}, nil, http.StatusOK
}

func LoginUser(payload requests.Login) (response responses.AuthResponse, err error, status int) {
	var user models.User

	_ = database.Mysql.Where("email = ?", payload.Email).Find(&user)

	if user.Empty() {
		return response, errors.New("invalid credentials"), http.StatusUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return responses.AuthResponse{}, errors.New("invalid credentials"), http.StatusUnauthorized
	}

	tokenStore := uuid.New().String()

	tokenExpiry := time.Hour * (24 * 7)

	err = database.Redis.Set(
		context.Background(),
		"user_auth_"+tokenStore,
		user.Id,
		tokenExpiry,
	).Err()

	if err != nil {
		return responses.AuthResponse{}, helpers.ServerError(err), http.StatusInternalServerError
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(tokenExpiry).Unix(),
		"token": tokenStore,
	})

	token, err := at.SignedString([]byte(config.AppConfig.AppKey))

	if err != nil {
		return responses.AuthResponse{}, helpers.ServerError(err), http.StatusInternalServerError
	}

	if !user.EmailVerified() {
		otp, err := gotp.New(gotp.Config{Redis: database.Redis})

		if err != nil {
			return response, helpers.ServerError(err), http.StatusInternalServerError
		}

		otpToken, err := otp.Generate(gotp.Generate{
			Format:     gotp.NUMERIC,
			Length:     6,
			Identifier: fmt.Sprintf("signup_%s", user.Id),
			Expires:    time.Minute * 10,
		})

		if err != nil {
			return response, helpers.ServerError(err), http.StatusInternalServerError
		}

		err = mailer.EnqueueEmailTask(queue.Client, mailer.EmailPayload{
			TemplateName: "signup_otp",
			To:           user.Email,
			Subject:      "Signup OTP",
			Data: map[string]interface{}{
				"Token": otpToken,
				"Name":  user.Name,
				"Email": user.Email,
			},
		})

		if err != nil {
			return response, helpers.ServerError(err), http.StatusInternalServerError
		}
	}

	return responses.AuthResponse{
		Token: token,
		User:  responses.GenerateUserResponse(user),
	}, nil, http.StatusOK
}

func VerifyUser(r *http.Request, payload requests.VerifyUser) (message map[string]string, err error, status int) {
	user := middleware.GetUser(r.Context())

	otp, err := gotp.New(gotp.Config{Redis: database.Redis})

	if err != nil {
		return nil, helpers.ServerError(err), http.StatusInternalServerError
	}

	valid, err := otp.Verify(gotp.Verify{
		Identifier: fmt.Sprintf("signup_%s", user.Id),
		Token:      payload.Token,
	})

	if err != nil || !valid {
		return nil, errors.New("invalid token"), http.StatusNotAcceptable
	}

	err = database.Mysql.Where("id = ?", user.Id).
		Model(models.User{}).
		Update("email_verified_at", time.Now()).
		Error

	if err != nil {
		return nil, helpers.ServerError(err), http.StatusInternalServerError
	}

	return helpers.PrepareMessage("token verified"), nil, http.StatusOK
}

func PreForgot(r *http.Request, payload requests.PreForgot) (message map[string]string, err error, status int) {
	var user models.User

	err = database.Mysql.Where("email = ?", payload.Email).First(&user).Error

	if err != nil {
		return nil, helpers.ServerError(err), http.StatusBadRequest
	}

	if !user.Empty() {
		token := uuid.New().String()

		err = database.Redis.Set(r.Context(), "forgot_password_"+token, user.Id, time.Hour*3).Err()

		if err != nil {
			return nil, helpers.ServerError(err), http.StatusBadRequest
		}

		err = mailer.EnqueueEmailTask(queue.Client, mailer.EmailPayload{
			TemplateName: "forgot_password",
			To:           user.Email,
			Subject:      "Password reset",
			Data: map[string]interface{}{
				"Token": token,
				"Email": payload.Email,
			},
		})

		if err != nil {
			return nil, helpers.ServerError(err), http.StatusInternalServerError
		}
	}

	return helpers.PrepareMessage("check your email"), nil, http.StatusOK
}

func PostForgot(r *http.Request, payload requests.PostForgot) (message map[string]string, err error, status int) {
	var user models.User

	userId, err := database.Redis.Get(r.Context(), "forgot_password_"+payload.Token).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errors.New("invalid token"), http.StatusNotAcceptable
		}

		return nil, helpers.ServerError(err), http.StatusInternalServerError
	}

	err = database.Mysql.Where("id = ?", userId).First(&user).Error

	if err != nil {
		return nil, helpers.ServerError(err), http.StatusInternalServerError
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, helpers.ServerError(err), http.StatusInternalServerError
	}

	err = database.Mysql.Model(&user).Updates(map[string]interface{}{
		"password": hashedPassword,
	}).Error

	return helpers.PrepareMessage("password reset complete"), nil, http.StatusOK
}
