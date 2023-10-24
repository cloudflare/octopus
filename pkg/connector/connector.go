//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package connector

import (
	"time"

	"github.com/cloudflare/octopus/pkg/model"
)

// The Connector (think Tentacle of the Octopus) is the glue between any given data source and the Octopus.
// It is responsible for getting the relevant data out of the data source and caching it internally for resilience.
// The Octopus will periodically ask all the connectors to enrich a new topology with the current data set to for them full enriched topology.
type Connector interface {
	GetName() string                      // Who am I?
	InitialLoad() error                   // Initial load of data
	Healthy() bool                        // Is the data source healthy?
	StartRefreshRoutine()                 // Start background update task
	EnrichTopology(*model.Topology) error // Update the given Topology with information from Connector
	GetLoadDuration() time.Duration       // How long did the last data load take?
	GetLoadTime() time.Time               // When was the current connector data loaded?
	GetUpdateErrorCount() uint64          // The number of time the refresh of connector data has failed
}
