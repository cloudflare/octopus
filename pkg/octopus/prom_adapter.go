//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package octopus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	topologyBuildDuration    = prometheus.NewDesc("octopus_topology_update_duration", "Time it took to build the topology (milliseconds)", nil, nil)
	topologyBuildTime        = prometheus.NewDesc("octopus_topology_build_time", "Timestamp (epoch) when the current topology was build", nil, nil)
	topologyItemCount        = prometheus.NewDesc("octopus_topology_item_count", "The number of instances per item", []string{"item_type"}, nil)
	connectorHealthyVec      = prometheus.NewDesc("octopus_connector_health", "Connector health indicatior (0/1)", []string{"connector"}, nil)
	connectorLoadDurationVec = prometheus.NewDesc("octopus_connector_load_duraton", "Timestamp (epoch) when the current connector data was fetched", []string{"connector"}, nil)
	connectorLoadTimeVec     = prometheus.NewDesc("octopus_connector_load_time", "Time it took to fetch data (milliseconds)", []string{"connector"}, nil)
	connectorUpdateErrorVec  = prometheus.NewDesc("octopus_connector_update_error_count", "The number of time the refresh of connector data has failed", []string{"connector"}, nil)
)

type PromAdapter struct {
	octopus *Octopus
}

func NewPromAdapter(octopus *Octopus) *PromAdapter {
	return &PromAdapter{
		octopus: octopus,
	}
}

func (p *PromAdapter) Describe(ch chan<- *prometheus.Desc) {
	ch <- topologyBuildDuration
	ch <- topologyBuildTime
	ch <- topologyItemCount
	ch <- connectorHealthyVec
	ch <- connectorLoadDurationVec
	ch <- connectorLoadTimeVec
	ch <- connectorUpdateErrorVec
}

func (p *PromAdapter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(topologyBuildDuration, prometheus.GaugeValue, float64(p.octopus.topologyBuildDuration.Load()))
	ch <- prometheus.MustNewConstMetric(topologyBuildTime, prometheus.GaugeValue, float64(p.octopus.topologyBuildTime.Load()))

	t := p.octopus.GetTopology()
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Sites)), "sites")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Pops)), "pops")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Colos)), "colos")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Nodes)), "devices")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Interfaces)), "interfaces")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Cables)), "cables")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Circuits)), "circuits")
	ch <- prometheus.MustNewConstMetric(topologyItemCount, prometheus.GaugeValue, float64(len(t.Prefixes)), "prefixes")

	for _, c := range p.octopus.connectors {
		ch <- prometheus.MustNewConstMetric(connectorHealthyVec, prometheus.GaugeValue, healthyToFloat64(c.Healthy()), c.GetName())
		ch <- prometheus.MustNewConstMetric(connectorLoadDurationVec, prometheus.GaugeValue, float64(c.GetLoadDuration().Milliseconds()), c.GetName())
		ch <- prometheus.MustNewConstMetric(connectorLoadTimeVec, prometheus.GaugeValue, float64(c.GetLoadTime().Unix()), c.GetName())
		ch <- prometheus.MustNewConstMetric(connectorUpdateErrorVec, prometheus.CounterValue, float64(c.GetUpdateErrorCount()), c.GetName())
	}
}

func healthyToFloat64(health bool) float64 {
	if health {
		return 1
	}

	return 0
}
