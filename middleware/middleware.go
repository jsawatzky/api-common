package middleware

import (
	"github.com/gorilla/mux"
	"github.com/jsawatzky/api-common/logging"
	"github.com/jsawatzky/api-common/metrics"
	"github.com/jsawatzky/api-common/recovery"
)

func Apply(r *mux.Router) {
	r.Use(recovery.Middleware)
	r.Use(logging.Middleware)
	r.Use(metrics.Middleware)
}
