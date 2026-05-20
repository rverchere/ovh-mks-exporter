package main

import (
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rverchere/ovh-mks-exporter/internal"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.Info("Starting application...")
	// https://www.ovh.com/auth/api/createToken
	client, err := ovh.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	exporter := internal.Exporter{
		Client:      client,
		ServiceName: os.Getenv("OVH_CLOUD_PROJECT_SERVICE"),
	}
	if err := exporter.NewExporter(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
