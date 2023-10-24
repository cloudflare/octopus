//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package octopus

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cloudflare/octopus/pkg/connector"
	"github.com/cloudflare/octopus/pkg/model"
	octopuspb "github.com/cloudflare/octopus/proto/octopus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	log "github.com/sirupsen/logrus"
)

const topologyRefreshTime = time.Minute

type Octopus struct {
	grpcPort uint16

	connectors            []connector.Connector
	topology              *model.Topology
	topologyMu            sync.RWMutex
	topologyBuildDuration atomic.Int64
	topologyBuildTime     atomic.Int64
}

// NewOctopus creates a new Octopus
func NewOctopus(grpcPort uint16) *Octopus {
	return &Octopus{
		grpcPort:   grpcPort,
		connectors: make([]connector.Connector, 0),
	}
}

// Init initializes the Ocotopus with the given list of connectors and triggers and initial load of data into the connectors
func (o *Octopus) Init(connectors []connector.Connector) error {
	for _, c := range connectors {
		log.Infof("Doing initial load for Connector %s...", c.GetName())
		err := c.InitialLoad()
		if err != nil {
			return fmt.Errorf("Initial load for Connector %s failed: %v\n", c.GetName(), err)
		}
	}

	o.connectors = connectors
	return nil
}

// Start starst the Octopus and connectors update routines as well as http + gRPC server
func (o *Octopus) Start() {
	for _, c := range o.connectors {
		c.StartRefreshRoutine()
	}

	go o.topologyRefreshRoutine()
	go o.serveGrpc()
}

// UpdateTopology triggers and instant update of the topology data from all configured connectors
func (o *Octopus) UpdateTopology() error {
	// Build new Topology
	topology := model.NewTopology()

	log.Info("Building new topology...")
	startTime := time.Now()

	for _, c := range o.connectors {
		if !c.Healthy() {
			return fmt.Errorf("Connector %s is not healthy, not updating topology!", c.GetName())
		}

		log.Infof("Enriching topology with data from Connector %s...", c.GetName())
		err := c.EnrichTopology(topology)
		if err != nil {
			return fmt.Errorf("Enriching topology with data from Connector %s failed: %v", c.GetName(), err)
		}
	}

	// We got ourselves a new topology, add the time when we built it and store it
	topology.Timestamp = time.Now()
	o.topologyBuildDuration.Store(topology.Timestamp.Sub(startTime).Milliseconds())
	o.topologyBuildTime.Store(topology.Timestamp.Unix())

	o.topologyMu.Lock()
	defer o.topologyMu.Unlock()
	o.topology = topology

	return nil
}

// GetTopology returns a pointer to the current topology
func (o *Octopus) GetTopology() *model.Topology {
	o.topologyMu.RLock()
	defer o.topologyMu.RUnlock()

	return o.topology
}

func (o *Octopus) Healthy() bool {
	o.topologyMu.RLock()
	defer o.topologyMu.RUnlock()

	return o.topology != nil
}

func (o *Octopus) topologyRefreshRoutine() {
	for {
		<-time.After(topologyRefreshTime)

		err := o.UpdateTopology()
		if err != nil {
			log.Errorf("Failed to update topology: %v", err)
		}
	}
}

func (o *Octopus) serveGrpc() {
	portStr := fmt.Sprintf(":%d", o.grpcPort)
	log.Infof("Starting gRPC API server at %s", portStr)

	os := newOctopusServer(o)
	s := grpc.NewServer()
	s.RegisterService(&octopuspb.OctopusService_ServiceDesc, os)

	// Allow client to retrieve proto definition
	reflection.Register(s)

	l, err := net.Listen("tcp", portStr)
	if err != nil {
		log.Fatalf("Failed to listen on TCP port %d", o.grpcPort)
	}

	err = s.Serve(l)
	if err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
