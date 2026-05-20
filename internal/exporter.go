package internal

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ovh/go-ovh/ovh"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	Client      *ovh.Client
	ServiceName string
	MaxRetries  int
	S3Regions   []string
	LBRegions   []string
}

func (exporter *Exporter) NewExporter() error {
	err := prometheus.Register(&collector{exporter: exporter})
	if err != nil {
		if registered, ok := err.(prometheus.AlreadyRegisteredError); ok {
			prometheus.Unregister(registered.ExistingCollector)
			prometheus.MustRegister(&collector{exporter: exporter})
		} else {
			return err
		}
	}

	log.Info("Starting exporter, enjoy!")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:    ":9101",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Info("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server shutdown error: ", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
