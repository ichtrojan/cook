package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/ichtrojan/cook/helpers"
)

func VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUser(r.Context()).EmailVerifiedAt == nil {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusForbidden)

			_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("email not verified"))

			return
		}

		next.ServeHTTP(w, r)
	})
}
