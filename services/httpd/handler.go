package httpd

import (
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	// import static fs files
)

const (
	// DefaultChunkSize specifies the maximum number of points that will
	// be read before sending results back to the engine.
	//
	// This has no relation to the number of bytes that are returned.
	DefaultChunkSize = 10000

	// MaxFileSize specifies a maxium upload file size
	MaxFileSize = 1 * 1000 * 1000 * 1000 * 1000
)

// AuthenticationMethod defines the type of authentication used.
type AuthenticationMethod int

// Supported authentication methods.
const (
	UserAuthentication AuthenticationMethod = iota
	BearerAuthentication
)

// TODO: Check HTTP response codes: 400, 401, 403, 409.

// Route specifies how to handle a HTTP verb for a given endpoint.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(http.ResponseWriter, *http.Request)
}

// Handler represents an HTTP handler for the InfluxDB server.
type Handler struct {
	mux     *pat.PatternServeMux
	Version string

	Kafka  *kafka.Service
	Config *Config
	Logger *logrus.Entry
	Close  chan struct{}
}

// NewHandler returns a new instance of handler with routes.
func NewHandler(c Config) *Handler {
	h := &Handler{
		mux:    pat.New(),
		Config: &c,
		Close:  make(chan struct{}),
	}
	h.AddRoutes([]Route{
		Route{
			"health-check", // Return a health check
			"GET", "/healthz", h.serveHealthCheck,
		},
		Route{ // Ping
			"ping-head",
			"HEAD", "/healthz", h.serveHealthCheck,
		},
		Route{ // Ping
			"metrics",
			"GET", "/metrics", promhttp.Handler().ServeHTTP,
		},
		Route{
			"mail",
			"POST", "/mail", h.acceptInboundEmail,
		},
	}...)

	return h
}

// AddRoutes sets the provided routes on the handler.
func (h *Handler) AddRoutes(routes ...Route) {
	for _, r := range routes {
		h.mux.Add(r.Method, r.Pattern, http.HandlerFunc(r.HandlerFunc))
	}

}

// ServeHTTP responds to HTTP request to the handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Cloudive-Version", h.Version)

	if strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		switch r.URL.Path {
		case "/debug/pprof/cmdline":
			pprof.Cmdline(w, r)
		case "/debug/pprof/profile":
			pprof.Profile(w, r)
		case "/debug/pprof/symbol":
			pprof.Symbol(w, r)
		default:
			pprof.Index(w, r)
		}
	} else {
		w.Header().Add("Content-Type", "application/json")
		h.mux.ServeHTTP(w, r)
	}

	// atomic.AddInt64(&h.stats.RequestDuration, time.Since(start).Nanoseconds())
}

// serveHealthCheck returns an empty response to comply with OPTIONS pre-flight requests
func (h *Handler) serveHealthCheck(w http.ResponseWriter, r *http.Request) {
	h.writeHeader(w, http.StatusNoContent)
}

// serveOptions returns an empty response to comply with OPTIONS pre-flight requests
func (h *Handler) serveOptions(w http.ResponseWriter, r *http.Request) {
	h.writeHeader(w, http.StatusNoContent)
}

// writeHeader writes the provided status code in the response, and
// updates relevant http error statistics.
func (h *Handler) writeHeader(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
