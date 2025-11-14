package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func (sr *statusRecorder) Write(b []byte) (int, error) {
	if sr.status == 0 {
		sr.status = http.StatusOK
	}
	n, err := sr.ResponseWriter.Write(b)
	sr.bytes += n
	return n, err
}

func Logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sr := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(sr, r)

			lat := time.Since(start)
			logger.Printf(`req_id=%s method=%s path=%s status=%d bytes=%d latency_ms=%.3f`,
				GetRequestID(r.Context()), r.Method, r.URL.Path, sr.status, sr.bytes, float64(lat.Microseconds())/1000.0)
		})
	}
}
