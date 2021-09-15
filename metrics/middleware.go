package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jsawatzky/api-common/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	httpRequestsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_size_bytes",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "code"},
	)
)

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path := "unknown"
		if route != nil {
			if pt, err := route.GetPathTemplate(); err == nil {
				path = pt
			}
		}

		start := time.Now()

		rec := internal.RecordResponse(rw)
		h.ServeHTTP(rec, r)

		labels := []string{r.Method, path, strconv.Itoa(rec.Status())}
		httpRequestsDuration.WithLabelValues(labels...).Observe(float64(time.Since(start).Seconds()))
		httpResponseSize.WithLabelValues(labels...).Observe(float64(rec.ResponseSize()))
		httpRequestsTotal.WithLabelValues(labels...).Inc()
	})
}
