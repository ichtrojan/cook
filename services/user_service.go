package services

import (
	"github.com/ichtrojan/cook/database"
	"github.com/ichtrojan/cook/helpers"
	middleware "github.com/ichtrojan/cook/middlewares"
	"github.com/ichtrojan/cook/models"
	"github.com/ichtrojan/cook/responses"
	"net/http"
)

func GetUser(r *http.Request) (response responses.UserResponse, err error, status int) {
	var user models.User

	err = database.Mysql.Where("id = ?", middleware.GetUser(r.Context()).Id).
		Preload("Businesses.Members.User").
		Find(&user).Error

	if err != nil {
		return responses.UserResponse{}, helpers.ServerError(err), http.StatusInternalServerError
	}

	return responses.GenerateUserResponse(user), nil, http.StatusOK
}
