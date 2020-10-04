package usvc

import (
	"flag"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricOpts struct{}

func (o *MetricOpts) Flags(fs *flag.FlagSet) {}

func (o MetricOpts) Metrics(mux *http.ServeMux) (live, ready HealthProbe) {
	mux.Handle("/metrics", promhttp.Handler())

	live.Healthy(true)
	ready.Healthy(true)
	mux.Handle("/liveness", live)
	mux.Handle("/readiness", ready)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	return
}

type HealthProbe bool

func (h *HealthProbe) Healthy(b bool) {
	*h = HealthProbe(b)
}

func (h HealthProbe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
