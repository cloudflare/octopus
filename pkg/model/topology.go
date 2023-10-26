//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	"fmt"
	"sort"
	"time"

	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

type Topology struct {
	Timestamp time.Time

	Sites                map[string]*Site
	Pops                 map[string]*Pop
	Colos                map[uint16]*Colo
	Nodes                map[string]*Device
	Interfaces           map[int64]*Interface
	DevicesByInterfaceID map[int64]*Device
	Cables               map[string]*Cable
	Prefixes             map[int64]*Prefix
	Circuits             map[string]*Circuit
}

func NewTopology() *Topology {
	return &Topology{
		Sites:                make(map[string]*Site),
		Pops:                 make(map[string]*Pop),
		Colos:                make(map[uint16]*Colo),
		Nodes:                make(map[string]*Device),
		Interfaces:           make(map[int64]*Interface),
		DevicesByInterfaceID: make(map[int64]*Device),
		Cables:               make(map[string]*Cable),
		Prefixes:             make(map[int64]*Prefix),
		Circuits:             make(map[string]*Circuit),
	}
}

// GetAddDevice will return the pointer to the Device with the given name.
// If no device with that name currently exists, it will be created.
func (t *Topology) AddDeviceIfNotExists(name string) *Device {
	dev, exists := t.Nodes[name]
	if !exists {
		dev = NewDevice(name)
		t.Nodes[name] = dev
	}

	return dev
}

func (t *Topology) GetDevice(name string) *Device {
	return t.Nodes[name]
}

func (t *Topology) AddSiteIfNotExists(name string) *Site {
	site, exists := t.Sites[name]
	if !exists {
		site = NewSite(name)
		t.Sites[name] = site
	}

	return site
}

func (t *Topology) AddPopIfNotExists(name string) *Pop {
	pop, exists := t.Pops[name]
	if !exists {
		pop = NewPop(name)
		t.Pops[name] = pop
	}

	return pop
}

func (t *Topology) AddColoIfNotExists(id uint16, name string, popName string) *Colo {
	pop := t.AddPopIfNotExists(popName)

	colo, exists := t.Colos[id]
	if !exists {
		colo = NewColo(id, name, pop)
		t.Colos[id] = colo
		pop.Colos = append(pop.Colos, colo)
	}

	return colo
}

func (t *Topology) GetColo(id uint16) *Colo {
	return t.Colos[id]
}

func (t *Topology) ToProto() *octopuspb.Topology {
	if t == nil {
		return nil
	}

	protoTopology := &octopuspb.Topology{
		Timestamp: uint64(t.Timestamp.Unix()),
		Devices:   make([]*octopuspb.Device, 0),
	}

	for _, dev := range t.Nodes {
		protoTopology.Devices = append(protoTopology.Devices, dev.ToProto())
	}

	if len(t.Sites) > 0 {
		protoTopology.Sites = make([]*octopuspb.Site, 0)
		for _, site := range t.Sites {
			protoTopology.Sites = append(protoTopology.Sites, site.ToProto())
		}
	}

	if len(t.Pops) > 0 {
		protoTopology.Pops = make([]*octopuspb.Pop, 0)
		for _, pop := range t.Pops {
			protoTopology.Pops = append(protoTopology.Pops, pop.ToProto())
		}
	}

	if len(t.Colos) > 0 {
		protoTopology.Colos = make([]*octopuspb.Colo, 0)
		for _, colo := range t.Colos {
			protoTopology.Colos = append(protoTopology.Colos, colo.ToProto())
		}
	}

	if len(t.Cables) > 0 {
		protoTopology.Cables = make([]*octopuspb.Cable, 0)
		for _, cable := range t.Cables {
			protoTopology.Cables = append(protoTopology.Cables, cable.ToProto())
		}
	}

	if len(t.Circuits) > 0 {
		for _, ckt := range t.Circuits {
			protoTopology.Circuits = append(protoTopology.Circuits, ckt.ToProto())
		}
	}

	if len(t.Prefixes) > 0 {
		protoTopology.Prefixes = make([]*octopuspb.Prefix, 0)
		for _, prefix := range t.Prefixes {
			protoTopology.Prefixes = append(protoTopology.Prefixes, prefix.ToProto())
		}
	}

	sortTopology(protoTopology)
	return protoTopology
}

func sortTopology(topology *octopuspb.Topology) {
	sort.Slice(topology.Devices, func(i, j int) bool {
		return topology.Devices[i].Name < topology.Devices[j].Name
	})

	for _, d := range topology.Devices {
		sort.Slice(d.Interfaces, func(i, j int) bool {
			return d.Interfaces[i].Name < d.Interfaces[j].Name
		})

		for _, ifa := range d.Interfaces {
			sort.Slice(ifa.Units, func(i, j int) bool {
				if ifa.Units[0].OuterTag != ifa.Units[1].OuterTag {
					return ifa.Units[0].OuterTag > ifa.Units[1].OuterTag
				}

				return ifa.Units[0].InnerTag > ifa.Units[1].InnerTag
			})
		}

		sort.Slice(d.FrontPorts, func(i, j int) bool {
			return d.FrontPorts[i].Name < d.FrontPorts[j].Name
		})

		sort.Slice(d.RearPorts, func(i, j int) bool {
			return d.RearPorts[i].Name < d.RearPorts[j].Name
		})
	}

	sort.Slice(topology.Sites, func(i, j int) bool {
		return topology.Sites[i].Name < topology.Sites[j].Name
	})

	sort.Slice(topology.Pops, func(i, j int) bool {
		return topology.Pops[i].Name < topology.Pops[j].Name
	})

	sort.Slice(topology.Colos, func(i, j int) bool {
		return topology.Colos[i].Name < topology.Colos[j].Name
	})

	sort.Slice(topology.Prefixes, func(i, j int) bool {
		return comparePrefixes(topology.Prefixes[i], topology.Prefixes[j])
	})

	sort.Slice(topology.Cables, func(i, j int) bool {
		return cableToString(topology.Cables[i]) < cableToString(topology.Cables[j])
	})
}

func cableToString(c *octopuspb.Cable) string {
	return fmt.Sprintf("%s:%s<->%s:%s", c.AEnd.DeviceName, c.AEnd.EndpointName, c.BEnd.DeviceName, c.BEnd.EndpointName)
}

func (t *Topology) DeviceAndInterfaceExists(devName string, ifName string) error {
	dev := t.GetDevice(devName)
	if dev == nil {
		return fmt.Errorf("device %s not found", devName)
	}

	if dev.GetInterface(ifName) == nil {
		return fmt.Errorf("interface %s:%s not found", devName, ifName)
	}

	return nil
}

func (t *Topology) FindInterfaceUnitByMetaDataAndRole(key string, value string, role string) []*InterfaceUnit {
	res := make([]*InterfaceUnit, 0)
	for _, d := range t.Nodes {
		if d.Role != role {
			continue
		}

		for _, ifa := range d.Interfaces {
			for _, u := range ifa.Units {
				v, exists := u.MetaData.SemanticTags[key]
				if !exists {
					continue
				}

				if v == value {
					res = append(res, u)
				}
			}
		}
	}

	return res
}

// Return true if a < b
func comparePrefixes(a *octopuspb.Prefix, b *octopuspb.Prefix) bool {
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}

	if a.Prefix.Address.Higher > b.Prefix.Address.Higher {
		return false
	}

	if a.Prefix.Address.Higher < b.Prefix.Address.Higher {
		return true
	}

	if a.Prefix.Address.Lower > b.Prefix.Address.Lower {
		return false
	}

	if a.Prefix.Address.Lower < b.Prefix.Address.Lower {
		return true
	}

	return a.Prefix.Length < b.Prefix.Length
}
