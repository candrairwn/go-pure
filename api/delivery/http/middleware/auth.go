package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func IsAuthenticated(next http.Handler, log *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do something here

		next.ServeHTTP(w, r)
	})
}
