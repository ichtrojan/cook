package services

import (
	middleware "github.com/ichtrojan/cook/middlewares"
	"github.com/ichtrojan/cook/responses"
	"net/http"
)

func GetUser(r *http.Request) (response responses.UserResponse, err error, status int) {
	return responses.GenerateUserResponse(middleware.GetUser(r.Context())), nil, http.StatusOK
}
