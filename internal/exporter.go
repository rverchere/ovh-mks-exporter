package internal

import (
	"net/http"

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
}

// ListenAndServe : Convenience function to start exporter
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

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":9101", nil); err != nil {
		return err
	}
	return nil
}
