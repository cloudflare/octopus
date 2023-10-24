//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package netbox

import (
	"testing"
	"time"

	"github.com/cloudflare/octopus/pkg/model"
	"github.com/stretchr/testify/assert"

	bnet "github.com/bio-routing/bio-rd/net"
	dbModel "github.com/cloudflare/octopus/pkg/connector/netbox/model"
	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

func TestEnrichment(t *testing.T) {
	apiClient := &NetboxClient{
		db: &database{
			contentTypeDcimDevice:                 1,
			contentTypeDcimInterface:              2,
			contentTypeIpamIpaddress:              3,
			contentTypeIpamPrefix:                 4,
			contentTypeCircuitsCircuittermination: 5,
		},
	}

	tests := []struct {
		name     string
		nc       *NetboxConnector
		t        *model.Topology
		expected *octopuspb.Topology
		wantFail bool
	}{
		{
			name: "device only",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "ccr01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:       "ccr01.dus01",
						Role:       "ccr",
						SiteName:   "DUS01",
						Interfaces: []*octopuspb.Interface{},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
			},
		},
		{
			name: "failing test: device not existent",
			nc: &NetboxConnector{
				client:  apiClient,
				devices: map[int64]*dbModel.DcimDevice{},
				interfaces: map[int64]*dbModel.DcimInterface{
					1: {
						Name: "foo",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			wantFail: true,
		},
		{
			name: "failing test: interface not existent",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "ccr01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					1: {
						Name: "foo.100",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Parent: &dbModel.DcimInterface{
							Name: "foo",
						},
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			wantFail: true,
		},
		{
			name: "2 devices with interfaces and units",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "ccr01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
					2: {
						ID:   2,
						Name: "GCP",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CloudProvider",
							Slug: "cloud-provider",
						},
						Site: dbModel.DcimSite{
							Name: "ANY",
						},
						Tags: []string{
							"NET:ASN=16550",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					1: {
						Name: "Interconnect0",
						Device: dbModel.DcimDevice{
							Name: "GCP",
						},
						Type: "10gbase-x-sfpp",
						Tags: []string{
							"net:foo",
						},
					},
					2: {
						Name: "Ethernet0/0",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type:       "10gbase-x-sfpp",
						Tags:       []string{},
						MacAddress: "FE:FE:FE:DE:AD:FE",
						Speed:      10000000,
					},
					3: {
						Name: "Ethernet0/0.100",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "virtual",
						Parent: &dbModel.DcimInterface{
							Name: "Ethernet0/0",
						},
						Tags: []string{},
					},
					4: {
						Name: "Ethernet0/1",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "10gbase-x-sfpp",
						LAG: &dbModel.DcimInterface{
							Name: "bond0",
						},
						Tags: []string{},
					},
					5: {
						Name: "Ethernet0/2",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "10gbase-x-sfpp",
						LAG: &dbModel.DcimInterface{
							Name: "bond0",
						},
						Tags: []string{},
					},
					6: {
						Name: "bond0",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "lag",
						Tags: []string{},
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "ANY",
					},
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:     "GCP",
						Role:     "cloud-provider",
						SiteName: "ANY",
						Interfaces: []*octopuspb.Interface{
							{
								Name:  "Interconnect0",
								Type:  "10gbase-x-sfpp",
								Units: []*octopuspb.InterfaceUnit{},
								MetaData: &octopuspb.MetaData{
									SemanticTags: map[string]string{},
									Tags:         []string{"net:foo"},
								},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
						MetaData: &octopuspb.MetaData{
							SemanticTags: map[string]string{
								"NET:ASN": "16550",
							},
							Tags: []string{},
						},
					},
					{
						Name:     "ccr01.dus01",
						Role:     "ccr",
						SiteName: "DUS01",
						Interfaces: []*octopuspb.Interface{
							{
								Name: "Ethernet0/0",
								Type: "10gbase-x-sfpp",
								Units: []*octopuspb.InterfaceUnit{
									{
										Id:            100,
										InnerTag:      100,
										Ipv4Addresses: []*octopuspb.IPAddress{},
										Ipv6Addresses: []*octopuspb.IPAddress{},
									},
								},
							},
							{
								Name:        "Ethernet0/1",
								Type:        "10gbase-x-sfpp",
								LagMemberOf: "bond0",
								Units:       []*octopuspb.InterfaceUnit{},
							},
							{
								Name:        "Ethernet0/2",
								Type:        "10gbase-x-sfpp",
								LagMemberOf: "bond0",
								Units:       []*octopuspb.InterfaceUnit{},
							},
							{
								Name:  "bond0",
								Type:  "lag",
								Units: []*octopuspb.InterfaceUnit{},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
			},
		},
		{
			name: "2 devices with interfaces, units and IPs",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "ccr01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
					2: {
						ID:   2,
						Name: "GCP",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CloudProvider",
							Slug: "cloud-provider",
						},
						Site: dbModel.DcimSite{
							Name: "ANY",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					1: {
						ID:   1,
						Name: "Interconnect0",
						Device: dbModel.DcimDevice{
							Name: "GCP",
						},
						Type: "10gbase-x-sfpp",
						Tags: []string{
							"net:foo",
						},
					},
					2: {
						ID:   2,
						Name: "Ethernet0/0",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type:       "10gbase-x-sfpp",
						Tags:       []string{},
						MacAddress: "FE:FE:FE:DE:AD:FE",
						Speed:      10000000,
					},
					3: {
						ID:   3,
						Name: "Ethernet0/0.100",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "virtual",
						Parent: &dbModel.DcimInterface{
							Name: "Ethernet0/0",
						},
						Tags: []string{},
					},
					4: {
						ID:   4,
						Name: "Ethernet0/1",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "10gbase-x-sfpp",
						LAG: &dbModel.DcimInterface{
							Name: "bond0",
						},
						Tags: []string{},
					},
					5: {
						ID:   5,
						Name: "Ethernet0/2",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "10gbase-x-sfpp",
						LAG: &dbModel.DcimInterface{
							Name: "bond0",
						},
						Tags: []string{},
					},
					6: {
						ID:   6,
						Name: "bond0",
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Type: "lag",
						Tags: []string{},
					},
				},
				ipAddresses: []*dbModel.IpamIpaddress{
					{
						Address:              "192.0.2.0/31",
						AssignedObjectTypeID: 2,
						AssignedObjectID:     1,
					},
					{
						Address:              "192.0.2.1/31",
						AssignedObjectTypeID: 2,
						AssignedObjectID:     2,
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "ANY",
					},
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:     "GCP",
						Role:     "cloud-provider",
						SiteName: "ANY",
						Interfaces: []*octopuspb.Interface{
							{
								Name: "Interconnect0",
								Type: "10gbase-x-sfpp",
								Units: []*octopuspb.InterfaceUnit{
									{
										Id: 0,
										Ipv4Addresses: []*octopuspb.IPAddress{
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 0), 31).ToProto(),
											},
										},
										Ipv6Addresses: []*octopuspb.IPAddress{},
									},
								},
								MetaData: &octopuspb.MetaData{
									SemanticTags: map[string]string{},
									Tags:         []string{"net:foo"},
								},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
					{
						Name:     "ccr01.dus01",
						Role:     "ccr",
						SiteName: "DUS01",
						Interfaces: []*octopuspb.Interface{
							{
								Name: "Ethernet0/0",
								Type: "10gbase-x-sfpp",
								Units: []*octopuspb.InterfaceUnit{
									{
										Id: 0,
										Ipv4Addresses: []*octopuspb.IPAddress{
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 31).ToProto(),
											},
										},
										Ipv6Addresses: []*octopuspb.IPAddress{},
									},
									{
										Id:            100,
										InnerTag:      100,
										Ipv4Addresses: []*octopuspb.IPAddress{},
										Ipv6Addresses: []*octopuspb.IPAddress{},
									},
								},
							},
							{
								Name:        "Ethernet0/1",
								Type:        "10gbase-x-sfpp",
								LagMemberOf: "bond0",
								Units:       []*octopuspb.InterfaceUnit{},
							},
							{
								Name:        "Ethernet0/2",
								Type:        "10gbase-x-sfpp",
								LagMemberOf: "bond0",
								Units:       []*octopuspb.InterfaceUnit{},
							},
							{
								Name:  "bond0",
								Type:  "lag",
								Units: []*octopuspb.InterfaceUnit{},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
			},
		},
		{
			name: "cables + circuits only",
			nc: &NetboxConnector{
				client: apiClient,
				cables: []*dbModel.DcimCable{
					{
						TerminationATypeID: 2,
						TerminationBTypeID: 2,
						TerminationAID:     42,
						TerminationBID:     23,
					},
					{
						TerminationATypeID: 2,
						TerminationBTypeID: 5,
						TerminationAID:     98,
						TerminationBID:     1, // circuit termination
					},
					{
						TerminationATypeID: 2,
						TerminationBTypeID: 5,
						TerminationAID:     99,
						TerminationBID:     2, // circuit termination
					},

					// Cable with unknown termination type, e.g. serial or power cables, which are ignored for now
					{
						TerminationATypeID: 99,
						TerminationBTypeID: 5,
						TerminationAID:     99,
						TerminationBID:     2,
					},
					{
						TerminationATypeID: 1,
						TerminationBTypeID: 99,
						TerminationAID:     99,
						TerminationBID:     2,
					},
				},
				circuits: map[int64]*dbModel.CircuitsCircuit{
					1: {
						Cid:            "XCON-1234",
						TerminationAID: 1,
						TerminationZID: 2,
					},
				},
				circuitTerminations: map[int64]*dbModel.CircuitsCircuittermination{
					1: {
						ID:        1,
						CircuitID: 1,
					},
					2: {
						ID:        2,
						CircuitID: 1,
					},
				},
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "devA",
						Site: dbModel.DcimSite{
							Name: "SiteA",
						},
					},
					2: {
						ID:   2,
						Name: "devB",
						Site: dbModel.DcimSite{
							Name: "SiteA",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					42: {
						ID: 42,
						Device: dbModel.DcimDevice{
							Name: "devA",
						},
						Name: "ifaA",
					},
					23: {
						ID: 23,
						Device: dbModel.DcimDevice{
							Name: "devB",
						},
						Name: "ifaB",
					},
					98: {
						ID: 98,
						Device: dbModel.DcimDevice{
							Name: "devA",
						},
						Name: "ifaX",
					},
					99: {
						ID: 99,
						Device: dbModel.DcimDevice{
							Name: "devB",
						},
						Name: "ifaY",
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				Cables:               map[string]*model.Cable{},
				Circuits:             make(map[string]*model.Circuit),
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "SiteA",
					},
				},
				Pops:  make([]*octopuspb.Pop, 0),
				Colos: make([]*octopuspb.Colo, 0),
				Devices: []*octopuspb.Device{
					{
						Name:     "devA",
						SiteName: "SiteA",
						Interfaces: []*octopuspb.Interface{
							{
								Name:  "ifaA",
								Units: make([]*octopuspb.InterfaceUnit, 0),
							},
							{
								Name:  "ifaX",
								Units: make([]*octopuspb.InterfaceUnit, 0),
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
					{
						Name:     "devB",
						SiteName: "SiteA",
						Interfaces: []*octopuspb.Interface{
							{
								Name:  "ifaB",
								Units: make([]*octopuspb.InterfaceUnit, 0),
							},
							{
								Name:  "ifaY",
								Units: make([]*octopuspb.InterfaceUnit, 0),
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
				Cables: []*octopuspb.Cable{
					{
						AEnd: &octopuspb.CableEnd{
							DeviceName:   "devA",
							EndpointName: "ifaA",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_INTERFACE,
						},
						BEnd: &octopuspb.CableEnd{
							DeviceName:   "devB",
							EndpointName: "ifaB",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_INTERFACE,
						},
					},
					{
						AEnd: &octopuspb.CableEnd{
							DeviceName:   "devA",
							EndpointName: "ifaX",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_INTERFACE,
						},
						BEnd: &octopuspb.CableEnd{
							DeviceName:   "XCON-1234",
							EndpointName: "A",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_CIRCUIT_TERMINATION,
						},
					},
					{
						AEnd: &octopuspb.CableEnd{
							DeviceName:   "devB",
							EndpointName: "ifaY",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_INTERFACE,
						},
						BEnd: &octopuspb.CableEnd{
							DeviceName:   "XCON-1234",
							EndpointName: "Z",
							EndpointType: octopuspb.CableEndpointType_CABLE_ENDPOINT_TYPE_CIRCUIT_TERMINATION,
						},
					},
				},
				Circuits: []*octopuspb.Circuit{
					{
						Cid: "XCON-1234",
					},
				},
			},
		},
		{
			name: "test IPs on non units",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						Name: "ccr01.dus01",
						ID:   1,
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					1: {
						ID:       1,
						DeviceID: 1,
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Name: "lag1",
					},
				},
				ipAddresses: []*dbModel.IpamIpaddress{
					{
						Address:              "192.0.2.0/32",
						AssignedObjectID:     1,
						AssignedObjectTypeID: apiClient.db.contentTypeDcimInterface,
						CustomFieldData:      "{\"region_type\": \"sub-region\"}",
					},
					// Empty JSON object for CFD (as reported by current NetBox)
					{
						Address:              "192.0.2.1/32",
						AssignedObjectID:     1,
						AssignedObjectTypeID: apiClient.db.contentTypeDcimInterface,
						CustomFieldData:      "{}",
					},
					// Empty string for CFD
					{
						Address:              "192.0.2.2/32",
						AssignedObjectID:     1,
						AssignedObjectTypeID: apiClient.db.contentTypeDcimInterface,
						CustomFieldData:      "",
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:     "ccr01.dus01",
						Role:     "ccr",
						SiteName: "DUS01",
						Interfaces: []*octopuspb.Interface{
							{
								Name: "lag1",
								Units: []*octopuspb.InterfaceUnit{
									{
										Id: 0,
										Ipv4Addresses: []*octopuspb.IPAddress{
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 0), 32).ToProto(),
												MetaData: &octopuspb.MetaData{
													Tags:            make([]string, 0),
													SemanticTags:    make(map[string]string),
													CustomFieldData: "{\"region_type\": \"sub-region\"}",
												},
											},
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32).ToProto(),
											},
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 2), 32).ToProto(),
											},
										},
										Ipv6Addresses: []*octopuspb.IPAddress{},
									},
								},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
			},
		},
		{
			name: "patch panel",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "pp01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "PP",
							Slug: "pp",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
				},
				frontPorts: map[int64]*dbModel.DcimFrontport{
					1: {
						ID:               1,
						Name:             "FP-A",
						RearPortPosition: 1,
						DeviceID:         1,
						RearPortID:       2,
					},
				},
				rearPorts: map[int64]*dbModel.DcimRearport{
					2: {
						ID:        2,
						Name:      "RP-A",
						Positions: 4,
						DeviceID:  1,
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:       "pp01.dus01",
						Role:       "pp",
						SiteName:   "DUS01",
						Interfaces: []*octopuspb.Interface{},
						FrontPorts: []*octopuspb.FrontPort{
							{
								Name:             "FP-A",
								RearPort:         "RP-A",
								RearPortPosition: 1,
							},
						},
						RearPorts: []*octopuspb.RearPort{
							{
								Name:      "RP-A",
								Positions: 4,
							},
						},
					},
				},
			},
		},
		{
			name: "device with q-in-q interface",
			nc: &NetboxConnector{
				client: apiClient,
				devices: map[int64]*dbModel.DcimDevice{
					1: {
						ID:   1,
						Name: "ccr01.dus01",
						DeviceRole: dbModel.DcimDevicerole{
							Name: "CCR",
							Slug: "ccr",
						},
						Site: dbModel.DcimSite{
							Name: "DUS01",
						},
					},
				},
				interfaces: map[int64]*dbModel.DcimInterface{
					100: {
						ID:       100,
						DeviceID: 1,
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Name: "et-0/0/0",
					},
					222: {
						ID:       222,
						DeviceID: 1,
						Device: dbModel.DcimDevice{
							Name: "ccr01.dus01",
						},
						Name:     "et-0/0/0.23.42",
						Type:     "virtual",
						ParentID: 100,
						Parent: &dbModel.DcimInterface{
							ID:   100,
							Name: "et-0/0/0",
						},
					},
				},
				ipAddresses: []*dbModel.IpamIpaddress{
					{
						ID:                   4242,
						Address:              "169.254.0.0/31",
						AssignedObjectID:     222,
						AssignedObjectTypeID: apiClient.db.contentTypeDcimInterface,
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
			},
			expected: &octopuspb.Topology{
				Sites: []*octopuspb.Site{
					{
						Name: "DUS01",
					},
				},
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices: []*octopuspb.Device{
					{
						Name:     "ccr01.dus01",
						Role:     "ccr",
						SiteName: "DUS01",
						Interfaces: []*octopuspb.Interface{
							{
								Name: "et-0/0/0",
								Units: []*octopuspb.InterfaceUnit{
									{
										Id:       42,
										OuterTag: 23,
										InnerTag: 42,
										Ipv4Addresses: []*octopuspb.IPAddress{
											{
												IP: bnet.NewPfx(bnet.IPv4FromOctets(169, 254, 0, 0), 31).Ptr().ToProto(),
											},
										},
										Ipv6Addresses: make([]*octopuspb.IPAddress, 0),
									},
								},
							},
						},
						FrontPorts: []*octopuspb.FrontPort{},
						RearPorts:  []*octopuspb.RearPort{},
					},
				},
			},
		},
		{
			name: "prefixes",
			nc: &NetboxConnector{
				client: apiClient,
				prefixes: []*dbModel.IpamPrefix{
					{
						ID:     1,
						Prefix: "100.64.0.0/26",
						Tags: []string{
							"foo:bar",
						},
					},
					{
						ID:     2,
						Prefix: "192.0.2.0/24",
						Tags: []string{
							"isDocumentation=true",
						},
					},
				},
			},
			t: &model.Topology{
				Timestamp:            time.Unix(0, 0),
				Nodes:                map[string]*model.Device{},
				Sites:                map[string]*model.Site{},
				DevicesByInterfaceID: make(map[int64]*model.Device),
				Interfaces:           make(map[int64]*model.Interface),
				Prefixes:             make(map[int64]*model.Prefix),
			},
			expected: &octopuspb.Topology{
				Sites:    make([]*octopuspb.Site, 0),
				Pops:     make([]*octopuspb.Pop, 0),
				Colos:    make([]*octopuspb.Colo, 0),
				Cables:   make([]*octopuspb.Cable, 0),
				Circuits: make([]*octopuspb.Circuit, 0),
				Devices:  make([]*octopuspb.Device, 0),
				Prefixes: []*octopuspb.Prefix{
					{
						Prefix: bnet.NewPfx(bnet.IPv4FromOctets(100, 64, 0, 0), 26).ToProto(),
						MetaData: &octopuspb.MetaData{
							SemanticTags: map[string]string{},
							Tags: []string{
								"foo:bar",
							},
						},
					},
					{
						Prefix: bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 0), 24).ToProto(),
						MetaData: &octopuspb.MetaData{
							SemanticTags: map[string]string{
								"isDocumentation": "true",
							},
							Tags: []string{},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		err := test.nc._enrichTopology(test.t)
		if err != nil && !test.wantFail {
			t.Errorf("unexpected failure for test %q: %v", test.name, err)
			continue
		}

		if err == nil && test.wantFail {
			t.Errorf("unexpected success for test %q", test.name)
			continue
		}

		if err != nil && test.wantFail {
			continue
		}

		assert.Equal(t, test.expected, test.t.ToProto(), test.name)
	}
}

func TestExtractInterfaceAndUnit(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedName    string
		expectedVLANTag model.VLANTag
		wantFail        bool
	}{
		{
			name:            "dot1q",
			input:           "vlan.900",
			expectedName:    "vlan",
			expectedVLANTag: model.NewVLANTag(0, 900),
		},
		{
			name:            "Q-in-Q",
			input:           "xe-0/0/0.100.200",
			expectedName:    "xe-0/0/0",
			expectedVLANTag: model.NewVLANTag(100, 200),
		},
		{
			name:            "broken",
			input:           "xe-0/0/0",
			expectedName:    "xe-0/0/0",
			expectedVLANTag: model.NewVLANTag(100, 200),
			wantFail:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ifName, vt, err := extractInterfaceAndUnit(test.input)
			if err != nil {
				if test.wantFail {
					return
				}

				t.Errorf("unexpected failure of test %q: %v", test.name, err)
				return
			}

			if test.wantFail {
				t.Errorf("unexpected success of test %q", test.name)
				return
			}

			assert.Equal(t, test.expectedName, ifName, test.name)
			assert.Equal(t, test.expectedVLANTag, vt, test.name)
		})
	}
}
