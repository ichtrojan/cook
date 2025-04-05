package responses

import (
	"github.com/ichtrojan/cook/helpers"
	"github.com/ichtrojan/cook/models"
)

type UserResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

type SimplifiedUserResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

func GenerateUserResponse(user models.User) UserResponse {
	return UserResponse{
		Id:            user.Id,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified(),
		CreatedAt:     helpers.JSONTime{Time: user.CreatedAt}.Json(),
		UpdatedAt:     helpers.JSONTime{Time: user.UpdatedAt}.Json(),
	}
}

func GenerateSimplifiedUserResponse(user models.User) SimplifiedUserResponse {
	return SimplifiedUserResponse{
		Id:            user.Id,
		Name:          user.Name,
		Email:         user.Email,
		EmailVerified: user.EmailVerified(),
		CreatedAt:     helpers.JSONTime{Time: user.CreatedAt}.Json(),
		UpdatedAt:     helpers.JSONTime{Time: user.UpdatedAt}.Json(),
	}
}

func GenerateSimplifiedUsersResponse(users []models.User) []SimplifiedUserResponse {
	var response []SimplifiedUserResponse

	for _, user := range users {
		response = append(response, GenerateSimplifiedUserResponse(user))
	}

	if len(response) == 0 {
		return []SimplifiedUserResponse{}
	}

	return response
}
