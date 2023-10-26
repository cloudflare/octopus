//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package model

import (
	octopuspb "github.com/cloudflare/octopus/proto/octopus"
)

type Interface struct {
	Name        string
	Type        string
	LAGMemberOf string
	Units       map[VLANTag]*InterfaceUnit
	MetaData    *MetaData
}

type VLANTag struct {
	OuterTag uint16
	InnerTag uint16
}

func NewVLANTag(outerTag, innerTag uint16) VLANTag {
	return VLANTag{
		OuterTag: outerTag,
		InnerTag: innerTag,
	}
}

type InterfaceUnit struct {
	VLANTag
	ID            uint32
	IPv4Addresses []IP
	IPv6Addresses []IP
	MetaData      *MetaData
}

func newInterface(name string) *Interface {
	return &Interface{
		Name:     name,
		Units:    make(map[VLANTag]*InterfaceUnit),
		MetaData: NewMetaData(),
	}
}

func newInterfaceUnit(vlan VLANTag) *InterfaceUnit {
	return &InterfaceUnit{
		ID:            uint32(vlan.InnerTag),
		VLANTag:       vlan,
		IPv4Addresses: make([]IP, 0),
		IPv6Addresses: make([]IP, 0),
		MetaData:      NewMetaData(),
	}
}

func (iface *Interface) AddUnitIfNotExists(vlanTag VLANTag) *InterfaceUnit {
	u, exists := iface.Units[vlanTag]
	if !exists {
		u = newInterfaceUnit(vlanTag)
		iface.Units[vlanTag] = u
	}

	return u
}

func (iface *Interface) AddIPAddressIfNotExists(vlanTag VLANTag, newIP IP) {
	u := iface.AddUnitIfNotExists(vlanTag)

	if newIP.Address.Addr().IsIPv4() {
		u.IPv4Addresses = appendIPIfNotExists(u.IPv4Addresses, newIP)
		return
	}

	u.IPv6Addresses = appendIPIfNotExists(u.IPv6Addresses, newIP)
}

func (iface *Interface) ToProto() *octopuspb.Interface {
	if iface == nil {
		return nil
	}

	protoIface := &octopuspb.Interface{
		Name:        iface.Name,
		Type:        iface.Type,
		LagMemberOf: iface.LAGMemberOf,
		MetaData:    iface.MetaData.ToProto(),
	}

	if len(iface.Units) > 0 {
		protoIface.Units = make([]*octopuspb.InterfaceUnit, 0)
		for _, unit := range iface.Units {
			protoIface.Units = append(protoIface.Units, unit.ToProto())
		}
	}

	return protoIface
}

func (unit *InterfaceUnit) ToProto() *octopuspb.InterfaceUnit {
	if unit == nil {
		return nil
	}

	protoUnit := &octopuspb.InterfaceUnit{
		Id:       unit.ID,
		MetaData: unit.MetaData.ToProto(),
		OuterTag: uint32(unit.OuterTag),
		InnerTag: uint32(unit.InnerTag),
	}

	if len(unit.IPv4Addresses) > 0 {
		protoUnit.Ipv4Addresses = make([]*octopuspb.IPAddress, 0)
		for _, IP := range unit.IPv4Addresses {
			protoUnit.Ipv4Addresses = append(protoUnit.Ipv4Addresses, IP.ToProto())
		}
	}

	if len(unit.IPv6Addresses) > 0 {
		protoUnit.Ipv6Addresses = make([]*octopuspb.IPAddress, 0)
		for _, IP := range unit.IPv6Addresses {
			protoUnit.Ipv6Addresses = append(protoUnit.Ipv6Addresses, IP.ToProto())
		}
	}

	return protoUnit
}
