package observability

import (
	"crypto/subtle"
	"net/http"
	"net/http/pprof"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewAdminServer(address, pprofToken string, metrics *Metrics) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	pprofHandler := func(handler http.HandlerFunc) http.Handler {
		return protectPProf(pprofToken, handler)
	}
	mux.Handle("/debug/pprof/", pprofHandler(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", pprofHandler(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", pprofHandler(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", pprofHandler(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", pprofHandler(pprof.Trace))
	for _, profile := range []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"} {
		mux.Handle("/debug/pprof/"+profile, protectPProf(pprofToken, pprof.Handler(profile)))
	}

	return &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func protectPProf(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			provided := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
