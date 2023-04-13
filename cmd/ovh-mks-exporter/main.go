package main

import (
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/rverchere/ovh-mks-exporter/internal"

	log "github.com/sirupsen/logrus"
)

func main() {
	internal.ServiceName = os.Getenv("OVH_CLOUDPROJECT_SERVICENAME")
	internal.KubeId = os.Getenv("OVH_CLOUDPROJECT_KUBEID")
	var err error

	// https://www.ovh.com/auth/api/createToken
	internal.Client, err = ovh.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	exporter := internal.Exporter{}
	err = exporter.NewExporter()
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
