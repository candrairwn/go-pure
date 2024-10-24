package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/candrairwn/go-pure/api/utils"
)

func Recovery(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := utils.ResponseRecorder{ResponseWriter: w}
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler { // Handle the abort gracefully
					return
				}

				stack := make([]byte, 1024)
				n := runtime.Stack(stack, true)

				log.ErrorContext(r.Context(), "panic!",
					slog.Any("error", err),
					slog.String("stack", string(stack[:n])),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("query", r.URL.RawQuery),
					slog.String("ip", r.RemoteAddr))

				if wr.Status == 0 { // response is not written yet
					http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(&wr, r)
	})
}
