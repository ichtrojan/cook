package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ichtrojan/cook/helpers"
	"github.com/ichtrojan/cook/requests"
	"github.com/ichtrojan/cook/services"
	"github.com/thedevsaddam/govalidator"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var request requests.Register

	rules := govalidator.MapData{
		"name":     []string{"required", "alpha_space"},
		"email":    []string{"required", "email"},
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

	response, err, status := services.RegisterUser(request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))
		return
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(response)

	return
}

func Login(w http.ResponseWriter, r *http.Request) {
	var request requests.Login

	rules := govalidator.MapData{
		"email":    []string{"required", "email"},
		"password": []string{"required", "regex:^.{8,}"},
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

	response, err, status := services.LoginUser(request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))

		return
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(response)

	return
}

func VerifyUser(w http.ResponseWriter, r *http.Request) {
	var request requests.VerifyUser

	rules := govalidator.MapData{
		"token": []string{"numeric", "required"},
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

	message, err, status := services.VerifyUser(r, request)

	if err != nil {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage(err.Error()))

		return
	}

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(message)

	return
}
