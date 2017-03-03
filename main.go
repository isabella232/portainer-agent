package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/rancher/portainer-agent/healthcheck"
	"github.com/rancher/portainer-agent/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "portainer-agent"
	app.Usage = "Start the Portainer config agent"
	app.Action = launch

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "cattle-url",
			Usage:  "URL for cattle API",
			EnvVar: "CATTLE_URL",
		},
		cli.StringFlag{
			Name:   "cattle-access-key",
			Usage:  "Cattle API Access Key",
			EnvVar: "CATTLE_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "cattle-secret-key",
			Usage:  "Cattle API Secret Key",
			EnvVar: "CATTLE_SECRET_KEY",
		},
		cli.IntFlag{
			Name:   "health-check-port",
			Value:  10240,
			Usage:  "Port to configure an HTTP health check listener on",
			EnvVar: "HEALTH_CHECK_PORT",
		},
		cli.StringFlag{
			Name:  "config-file",
			Value: "endpoints.json",
			Usage: "Location of endpoints config file",
		},
	}

	app.Run(os.Args)
}

func launch(c *cli.Context) {
	resultChan := make(chan error)

	go func() {
		resultChan <- server.Watch(c.GlobalString("config-file"),
			c.GlobalString("cattle-access-key"),
			c.GlobalString("cattle-secret-key"),
			c.GlobalString("cattle-url"))
	}()

	go func() {
		resultChan <- healthcheck.StartHealthCheck(c.GlobalInt("health-check-port"))
	}()

	log.Fatalf("Exiting: %v", <-resultChan)
}
