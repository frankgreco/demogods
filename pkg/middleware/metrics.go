package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/frankgreco/demogods/pkg/log"
	"github.com/frankgreco/demogods/pkg/metrics"
)

// Metrics is a middleware that will log and report metrics
// corresponding to the current request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := httptest.NewRecorder()
		metrics.RequestInFlightCount.Inc()
		next.ServeHTTP(rec, r)
		metrics.RequestInFlightCount.Dec()
		endTime := time.Now()

		log.WithContext(r.Context()).Info("request details",
			zap.String("http_request_remote_address", r.RemoteAddr),
			zap.String("http_request_method", r.Method),
			zap.String("http_request_url_path", r.URL.Path),
			zap.String("http_request_duration", fmt.Sprintf("%gms", getReqDuration(startTime, endTime))),
			zap.Int("http_response_status_code", rec.Code),
		)

		metrics.RequestLatency.WithLabelValues(strconv.Itoa(rec.Code), r.Method).Observe(getReqDuration(startTime, endTime))

		for k, v := range rec.Header() {
			w.Header()[k] = v
		}
		w.Write(rec.Body.Bytes())
	})
}

func getReqDuration(start, finish time.Time) float64 {
	return (float64(finish.UnixNano()-start.UnixNano()) / float64(time.Millisecond))
}
