package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
)

var requestDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name: "rest-server_request_duration",
	Help: "The HTTP request latencies in milliseconds.",
}, []string{"method", "route", "status"})

//nolint
func init() {
	prometheus.MustRegister(requestDuration)
}

func Monitoring(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		duration := float64(m.Duration / time.Millisecond)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		requestDuration.WithLabelValues(
			r.Method, routePattern, strconv.Itoa(m.Code),
		).Observe(duration)
	})
}
