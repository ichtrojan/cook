package services

import (
	"net/http"

	middleware "github.com/ichtrojan/cook/middlewares"
	"github.com/ichtrojan/cook/responses"
)

func GetUser(r *http.Request) (response responses.UserResponse, err error, status int) {
	return responses.GenerateUserResponse(middleware.GetUser(r.Context())), nil, http.StatusOK
}
