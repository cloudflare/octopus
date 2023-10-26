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

func TestAddInterfaceItNotExists(t *testing.T) {
	dev := NewDevice("foo")

	tests := []struct {
		name           string
		ifName         string
		expectedDevice *Device
	}{
		{
			name:   "One iface foo",
			ifName: "foo",
			expectedDevice: &Device{
				Name: "foo",
				Interfaces: map[string]*Interface{
					"foo": {
						Name:     "foo",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
				},
				FrontPorts: make(map[string]*FrontPort),
				RearPorts:  make(map[string]*RearPort),
				MetaData:   NewMetaData(),
			},
		},
		{
			name:   "Adding iface foo again",
			ifName: "foo",
			expectedDevice: &Device{
				Name: "foo",
				Interfaces: map[string]*Interface{
					"foo": {
						Name:     "foo",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
				},
				FrontPorts: make(map[string]*FrontPort),
				RearPorts:  make(map[string]*RearPort),
				MetaData:   NewMetaData(),
			},
		},
		{
			name:   "Adding iface bar",
			ifName: "bar",
			expectedDevice: &Device{
				Name: "foo",
				Interfaces: map[string]*Interface{
					"foo": {
						Name:     "foo",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
					"bar": {
						Name:     "bar",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
				},
				FrontPorts: make(map[string]*FrontPort),
				RearPorts:  make(map[string]*RearPort),
				MetaData:   NewMetaData(),
			},
		},
		{
			name:   "Adding iface foo again",
			ifName: "foo",
			expectedDevice: &Device{
				Name: "foo",
				Interfaces: map[string]*Interface{
					"foo": {
						Name:     "foo",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
					"bar": {
						Name:     "bar",
						Units:    map[VLANTag]*InterfaceUnit{},
						MetaData: NewMetaData(),
					},
				},
				FrontPorts: make(map[string]*FrontPort),
				RearPorts:  make(map[string]*RearPort),
				MetaData:   NewMetaData(),
			},
		},
	}

	for _, test := range tests {
		dev.AddInterfaceItNotExists(test.ifName)
		assert.Equal(t, test.expectedDevice, dev)
	}
}

func TestDeviceToProto(t *testing.T) {
	devWithIface := NewDevice("foo")
	iface := devWithIface.AddInterfaceItNotExists("bar")
	iface.AddIPAddressIfNotExists(NewVLANTag(0, 0), NewIP(bnet.NewPfx(bnet.IPv4FromOctets(192, 0, 2, 1), 32)))

	tests := []struct {
		name        string
		device      *Device
		protoDevice *octopuspb.Device
	}{
		{
			name:        "nil",
			device:      nil,
			protoDevice: nil,
		},
		{
			name:   "Empty device",
			device: NewDevice("foo"),
			protoDevice: &octopuspb.Device{
				Name: "foo",
			},
		},
		{
			name:   "Device with 1 Interface",
			device: devWithIface,
			protoDevice: &octopuspb.Device{
				Name: "foo",
				Interfaces: []*octopuspb.Interface{
					{
						Name: "bar",
						Units: []*octopuspb.InterfaceUnit{
							{
								Id: 0,
								Ipv4Addresses: []*octopuspb.IPAddress{
									{
										IP: &bapi.Prefix{
											Address: &bapi.IP{Lower: 3221225985},
											Length:  32,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.protoDevice, test.device.ToProto(), test.name)
	}
}
