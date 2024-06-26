package leader

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"github.com/supertylerc/sre-portfolio/scheduler/pkg/leader"
)

func Metrics() {
	reg := prometheus.NewRegistry()
	leader.RegisterMetrics(reg)

	addr := fmt.Sprintf(":%s", viper.Get("METRICS_PORT"))
	mux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	mux.Handle("/metrics", promHandler)

	srv := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: time.Second,
		Handler:           mux,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Unable to start HTTP Server", slog.String("err", err.Error()))
	}
}
