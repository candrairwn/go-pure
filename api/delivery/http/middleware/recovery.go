package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/candrairwn/go-pure/api/utils"
	"go.uber.org/zap"
)

func Recovery(next http.Handler, log *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := utils.ResponseRecorder{ResponseWriter: w}
		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler { // Handle the abort gracefully
					return
				}

				stack := make([]byte, 1024)
				n := runtime.Stack(stack, true)

				log.Errorw("panic!",
					"error", err,
					"stack", string(stack[:n]),
					"method", r.Method,
					"path", r.URL.Path,
					"query", r.URL.RawQuery,
					"ip", r.RemoteAddr)

				if wr.Status == 0 { // response is not written yet
					http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(&wr, r)
	})
}
