package middleware

import (
	"banner-service/internal/metrics"
	"github.com/fatih/color"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type MetricsResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (m *MetricsResponseWriter) WriteHeader(code int) {
	m.statusCode = code
	m.ResponseWriter.WriteHeader(code)
}

func (m *MetricsResponseWriter) Status() int {
	if m.statusCode == 0 {
		return http.StatusOK
	}
	return m.statusCode
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &MetricsResponseWriter{ResponseWriter: w}

		defer func(t time.Time) {
			metrics.HTTPLatency.With(prometheus.Labels{"method": buildResource(r)}).Observe(time.Since(t).Seconds())
			if ww.Status() >= 500 {
				metrics.HTTPRequestTotalFail.With(prometheus.Labels{"method": buildResource(r)}).Inc()
				metrics.HTTPErrorCount.With(prometheus.Labels{"method": buildResource(r)}).Inc()
			} else {
				metrics.HTTPRequestTotalSuccess.With(prometheus.Labels{"method": buildResource(r)}).Inc()
			}
		}(time.Now())

		next.ServeHTTP(ww, r)
	})
}

func LogRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color.Yellow("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
