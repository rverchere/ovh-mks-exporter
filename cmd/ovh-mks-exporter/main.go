package main

import (
	"os"
	"strconv"

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

	maxRetries := 3
	if v := os.Getenv("OVH_MAX_RETRIES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxRetries = n
		}
	}

	exporter := internal.Exporter{
		Client:      client,
		ServiceName: os.Getenv("OVH_CLOUD_PROJECT_SERVICE"),
		MaxRetries:  maxRetries,
	}
	if err := exporter.NewExporter(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
