package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rverchere/ovh-mks-exporter/internal"

	log "github.com/sirupsen/logrus"
)

func parseRegions(v string) []string {
	var regions []string
	for _, r := range strings.Split(v, ",") {
		if r := strings.TrimSpace(r); r != "" {
			regions = append(regions, r)
		}
	}
	return regions
}

func setLogLevel() {
	level := os.Getenv("OVH_LOG_LEVEL")
	if level == "" {
		return
	}
	parsed, err := log.ParseLevel(level)
	if err != nil {
		log.Warnf("Invalid OVH_LOG_LEVEL %q, using default (info)", level)
		return
	}
	log.SetLevel(parsed)
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	setLogLevel()

	log.Infof("Starting application (version %s)...", internal.Version)
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
		S3Regions:   parseRegions(os.Getenv("OVH_S3_REGIONS")),
		LBRegions:   parseRegions(os.Getenv("OVH_LB_REGIONS")),
	}
	if err := exporter.NewExporter(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
