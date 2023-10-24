package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/cloudflare/octopus/pkg/connector"
	"github.com/cloudflare/octopus/pkg/connector/netbox"
	"github.com/cloudflare/octopus/pkg/octopus"
)

const (
	httpServerTimeout            = time.Second * 60
	netboxPostgresPasswordOption = "NETBOX_DB_PASSWORD"
)

var (
	grpcPort       = flag.Uint("grpc-port", 2342, "GRPC API server port")
	httpPort       = flag.Uint("http-port", 8080, "HTTP server port (for metrics)")
	mockConnectors = flag.Bool("mock-connectors", false, "If set, connectors will be used with mock data")

	netboxDisable      = flag.Bool("netbox.disable", false, "Disable NetBox connector")
	netboxDBHost       = flag.String("netbox.db.host", "localhost", "Netbox's postgres DB host")
	netboxDBPort       = flag.Uint("netbox.db.port", 5432, "Netbox's postgres DB port")
	netboxDBUser       = flag.String("netbox.db.user", "netbox", "Netbox's postgres DB user")
	netboxDBName       = flag.String("netbox.db.db", "netbox", "Netbox's postgres DB name")
	netboxDBPassword   = flag.String("netbox.db.password", "", fmt.Sprintf("Netbox DB password (should be set as ENV %q", netboxPostgresPasswordOption))
	netboxDBTLS        = flag.Bool("netbox.db.tls", true, "Use TLS for the DB connection")
	netboxDBCaCertPath = flag.String("netbox.db.ca-cert-file-path", "", "Path to CA certificate PEM file")
	netboxDBLogQueries = flag.Bool("netbox.db.log-queries", false, "Log DB queries")
)

func getConnectors() []connector.Connector {
	conns := make([]connector.Connector, 0)

	if !*netboxDisable {
		if *netboxDBPassword == "" {
			log.Fatalf("%s is a mandatory parameter", netboxPostgresPasswordOption)
		}

		conns = append(conns, netbox.NewConnector(*netboxDBHost, *netboxDBPort, *netboxDBUser, *netboxDBPassword, *netboxDBName, *netboxDBTLS, *netboxDBCaCertPath, *netboxDBLogQueries))
	}

	return conns
}

func getMockConnectors() []connector.Connector {
	log.Info("Running with mock connectors!")

	conns := make([]connector.Connector, 0)

	if !*netboxDisable {
		log.Info("Mock connector for NetBox not implemented (yet).")
	}

	return conns
}

func loadEnvVars() {
	netboxDBPasswordEnv := os.Getenv(netboxPostgresPasswordOption)
	if netboxDBPasswordEnv != "" {
		netboxDBPassword = &netboxDBPasswordEnv
	}
}

func main() {
	loadEnvVars()
	flag.Parse()

	log.Infof("Octopus starting...")

	/*
	 * Set up Connectors
	 */
	var connectors []connector.Connector
	if *mockConnectors {
		connectors = getMockConnectors()
	} else {
		connectors = getConnectors()
	}

	if len(connectors) == 0 {
		log.Fatal("No connectors enabled, srsly?")
	}

	/*
	 * Set up the Octopus
	 */
	o := octopus.NewOctopus(uint16(*grpcPort))
	err := o.Init(connectors)
	if err != nil {
		log.Fatalf("Failed to initialize octopus: %v", err)
	}

	err = o.UpdateTopology()
	if err != nil {
		log.Fatalf("Failed to update topology data: %v", err)
	}

	/*
	 * Set up Prometheus adapter
	 */
	prometheus.MustRegister(octopus.NewPromAdapter(o))
	go serveHTTP(o)

	/*
	 * Start Octopus
	 *
	 * This starts the connector + Octopus update routines, the http (metrics) listener,
	 * as well as the gRPC API server
	 */
	o.Start()

	select {}
}

func serveHTTP(o *octopus.Octopus) {
	portStr := fmt.Sprintf(":%d", *httpPort)
	log.Infof("Starting http server at %s", portStr)

	// Set up HTTP listener
	m := http.NewServeMux()
	s := &http.Server{
		Addr:           portStr,
		ReadTimeout:    httpServerTimeout,
		WriteTimeout:   httpServerTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        m,
	}

	m.HandleFunc("/ready", func(rw http.ResponseWriter, req *http.Request) {
		if o.Healthy() {
			_, _ = rw.Write([]byte("OK"))
			return
		}

		rw.WriteHeader(http.StatusServiceUnavailable)
		_, _ = rw.Write([]byte("NOK"))
	})

	m.Handle("/metrics", promhttp.Handler())

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("http.ListenAndServe failed: %v", err)
	}
}
