package middleware

import (
	"net/http"
	"time"

	"github.com/candrairwn/go-pure/api/utils"
	"go.uber.org/zap"
)

func Accesslog(next http.Handler, log *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wr := utils.ResponseRecorder{ResponseWriter: w}

		next.ServeHTTP(&wr, r)

		log.Infow("accessed",
			"latency", time.Since(start).String(),
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"ip", r.RemoteAddr,
			"status", wr.Status,
			"bytes", wr.NumBytes,
		)
	})
}
