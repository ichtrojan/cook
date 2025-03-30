package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/ichtrojan/cook/helpers"
	"io/ioutil"
	"net/http"
)

func AcceptJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ValidateJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("invalid JSON format"))

			return
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		var jsonTest interface{}
		if len(body) > 0 && json.Unmarshal(body, &jsonTest) != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("invalid JSON format"))

			return
		}

		next.ServeHTTP(w, r)
	})
}
