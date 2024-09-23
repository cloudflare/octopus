# Octopus

The Octopus is part of our network automation pipeline. Its job is to aggregate data from different authoritative data sources into one **Enriched Topology** to provide a full picture to downstream systems.

The enriched topology is a graph of the entire network including all devices as well as connections like cables, circuits, etc.
The nodes of the graph represent the devices -- routers, switches, optical equipment, servers, etc. -- including their meta data like status, device type, platform, role, and any custom fields or tags applied in NetBox for example.
The edges within the graph are direct cables (DACs, copper/fiber patches), cable paths (e.g. through patch panels or optical equipment), circuits (e.g. dark fibers, wawe lengths, transports, etc.), or even could be serial connections.

## Connectors

Each connector (think tentacle) of the Octopus taps into one of our data sources and consumes the bits we are interested in.
It is responsible for querying the data, caching it locally, and updating the data in an interval meaningful for the data source and obtaining useful triggers, if any.

## Topology generation

The Octopus holds the global Topology.

To gather data from all Connectors it will pass a pointer to a (single) new Topology object into each Connector, which will add its insight into relevant parts of the Topology.
If devices, interfaces of devices, or other attributes are missing in the Topology, it is the Connectors responsible to add them.

## Open questions

Should we just regenerate the Topology on every run (time based, trigger based, or both?) or should each Connector know (and therefore have the responsibility to figure out) if it has new data since the last run, so the Octopus can query all Connectors and the Topology only needs to be updates if at least one Connector has need data?

# Observability

The Octopus exposes a number of metrics via an HTTP endpoint ready to be scraped by Prometheus.

 * `octopus_topology_update_duration` - Time it took to build the topology (milliseconds)
 * `octopus_topology_build_time` - Timestamp (epoch) when the current topology was build
 * `octopus_topology_item_count` - The number of instances per item (broken out bylabel `item_type`)
 * `octopus_connector_health` - Connector health indicatior (0/1) (broken out bylabel `connector`)
 * `octopus_connector_load_duraton` - Timestamp (epoch) when the current connector data was fetched (broken out by label `connector`)
 * `octopus_connector_load_time` - Time it took to fetch data (milliseconds) (broken out by label `connector`)
 * `octopus_connector_update_error_count` - The number of time the refresh of connector data has failed (broken out by label `connector`)

 Other than those, Octopus is exposing gRPC-related metrics that comes from [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/providers/prometheus).

# Querying data

The Octopus exposes a gRPC API to query the enriched topology data.

You can manually query the enriched topology using `grpcurl` from the gRPC endpoint. Be aware that you need to increase the message size if the topology is larger than the default of 4MB.
A call could look like this, querying the `bond0` interface if `ccr01.pad01` 

```bash
grpcurl -max-msg-sz=100000000 octopus-production.example.com:443 cloudflare.net.octopus.OctopusService.GetTopology | jq '.topology.devices[] | select(.name=="ccr01.pad01") | .interfaces[] | select(.name=="bond0")'
```