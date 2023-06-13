package internal

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	ServiceName string
	KubeId      string
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
	log.Fatal(http.ListenAndServe(":9101", nil))

	if err != nil {
		log.Fatal("failed to start server: ", err)
	}

	return nil
}
