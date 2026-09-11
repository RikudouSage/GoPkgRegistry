package middleware

import (
	"GoPkgRepository/helper"
	"net/http"
	"strings"
)

func RequiresAuthorizationMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
				authHeader = after
			}

			if authHeader != apiKey {
				writer.WriteHeader(http.StatusUnauthorized)
				_ = helper.WriteJSON(writer, map[string]string{
					"error": "Unauthorized",
				})
				return
			}

			next.ServeHTTP(writer, req)
		})
	}
}
