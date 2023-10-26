//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	"testing"

	bnet "github.com/bio-routing/bio-rd/net"
	bapi "github.com/bio-routing/bio-rd/net/api"
	octopuspb "github.com/cloudflare/octopus/proto/octopus"
	"github.com/stretchr/testify/assert"
)

var ipv4IP_1 = bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32)
var ipv4IP_2 = bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 2), 32)
var ipv6IP_1 = bnet.NewPfx(bnet.IPv6FromBlocks(0x2001, 0xdb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x42), 128)

func TestAddIPAddressIfNotExists(t *testing.T) {
	iface := newInterface("foo")

	tests := []struct {
		name          string
		ip            bnet.Prefix
		unit          uint16
		expectedIface *Interface
	}{
		{
			name: "Add one IPv4 IP",
			unit: 0,
			ip:   ipv4IP_1,
			expectedIface: &Interface{
				Name: "foo",
				Units: map[VLANTag]*InterfaceUnit{
					NewVLANTag(0, 0): {
						IPv4Addresses: []IP{
							NewIP(ipv4IP_1),
						},
						IPv6Addresses: []IP{},
						MetaData:      NewMetaData(),
					},
				},
				MetaData: NewMetaData(),
			},
		},
		{
			name: "Add the same IPv4 IP again",
			ip:   ipv4IP_1,
			expectedIface: &Interface{
				Name: "foo",
				Units: map[VLANTag]*InterfaceUnit{
					NewVLANTag(0, 0): {
						IPv4Addresses: []IP{
							NewIP(ipv4IP_1),
						},
						IPv6Addresses: []IP{},
						MetaData:      NewMetaData(),
					},
				},
				MetaData: NewMetaData(),
			},
		},
		{
			name: "Add a 2nd IPv4 IP",
			ip:   ipv4IP_2,
			expectedIface: &Interface{
				Name: "foo",
				Units: map[VLANTag]*InterfaceUnit{
					NewVLANTag(0, 0): {
						IPv4Addresses: []IP{
							NewIP(ipv4IP_1),
							NewIP(ipv4IP_2),
						},
						IPv6Addresses: []IP{},
						MetaData:      NewMetaData(),
					},
				},
				MetaData: NewMetaData(),
			},
		},
		{
			name: "Add an IPv6 IP",
			ip:   ipv6IP_1,
			expectedIface: &Interface{
				Name: "foo",
				Units: map[VLANTag]*InterfaceUnit{
					NewVLANTag(0, 0): {
						IPv4Addresses: []IP{
							NewIP(ipv4IP_1),
							NewIP(ipv4IP_2),
						},
						IPv6Addresses: []IP{
							NewIP(ipv6IP_1),
						},
						MetaData: NewMetaData(),
					},
				},
				MetaData: NewMetaData(),
			},
		},
	}

	for _, test := range tests {
		iface.AddIPAddressIfNotExists(NewVLANTag(0, 0), NewIP(test.ip))
		assert.Equal(t, test.expectedIface, iface)
	}
}

func TestInterfaceToProto(t *testing.T) {
	ifaceWithIPs := newInterface("bar")
	ifaceWithIPs.AddIPAddressIfNotExists(NewVLANTag(0, 0), NewIP(ipv4IP_1))
	ifaceWithIPs.AddIPAddressIfNotExists(NewVLANTag(0, 0), NewIP(ipv6IP_1))
	ifaceWithIPs.MetaData = &MetaData{
		Tags: []string{
			"regular-tag",
		},
		SemanticTags: map[string]string{
			"foo": "bar",
		},
		CustomFieldData: "{ \"foo\": \"baz\"}",
	}

	tests := []struct {
		name       string
		iface      *Interface
		protoIface *octopuspb.Interface
	}{
		{
			name:       "nil",
			iface:      nil,
			protoIface: nil,
		},
		{
			name:  "Empty interface",
			iface: newInterface("foo"),
			protoIface: &octopuspb.Interface{
				Name: "foo",
			},
		},

		{
			name:  "Interface with IPs and MetaData",
			iface: ifaceWithIPs,
			protoIface: &octopuspb.Interface{
				Name: "bar",
				Units: []*octopuspb.InterfaceUnit{
					{
						Id: 0,
						Ipv4Addresses: []*octopuspb.IPAddress{
							{
								IP: &bapi.Prefix{
									Address: &bapi.IP{
										Lower: 3221225985,
									},
									Length: 32,
								},
							},
						},
						Ipv6Addresses: []*octopuspb.IPAddress{
							{
								IP: &bapi.Prefix{
									Address: &bapi.IP{
										Higher:  2306139568115548160,
										Lower:   66,
										Version: 1,
									},
									Length: 128,
								},
							},
						},
					},
				},
				MetaData: &octopuspb.MetaData{
					Tags: []string{
						"regular-tag",
					},
					SemanticTags: map[string]string{
						"foo": "bar",
					},
					CustomFieldData: "{ \"foo\": \"baz\"}",
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoIface, test.iface.ToProto(), test.name)
	}
}
