package controllers

import (
	"encoding/json"
	"github.com/ichtrojan/cook/helpers"
	"github.com/ichtrojan/cook/requests"
	"github.com/ichtrojan/cook/responses"
	"github.com/ichtrojan/cook/services"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

type GetUserResponse struct {
	User responses.UserResponse `json:"user"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	response, err, status := services.GetUser(r)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))

		return
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(GetUserResponse{response})

	return
}

func PreForgot(w http.ResponseWriter, r *http.Request) {
	var request requests.PreForgot

	rules := govalidator.MapData{
		"email": []string{"email", "required"},
	}

	opts := govalidator.Options{
		Request: r,
		Rules:   rules,
		Data:    &request,
	}

	validationErrors := helpers.ValidateRequest(opts, "json")

	if len(validationErrors) != 0 {
		helpers.ReturnValidationErrors(w, validationErrors)
		return
	}

	message, err, status := services.PreForgot(r, request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(message)
	return
}

func PostForgot(w http.ResponseWriter, r *http.Request) {
	var request requests.PostForgot

	rules := govalidator.MapData{
		"token":    []string{"uuid", "required"},
		"password": []string{"regex:^.{8,}", "required"},
	}

	opts := govalidator.Options{
		Request: r,
		Rules:   rules,
		Data:    &request,
	}

	validationErrors := helpers.ValidateRequest(opts, "json")

	if len(validationErrors) != 0 {
		helpers.ReturnValidationErrors(w, validationErrors)
		return
	}

	message, err, status := services.PostForgot(r, request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(message)
	return
}
