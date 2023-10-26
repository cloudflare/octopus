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

type Device struct {
	Name       string
	Status     string
	Role       string
	Platform   string
	DeviceType string

	Colo *Colo
	Site *Site

	Interfaces map[string]*Interface
	FrontPorts map[string]*FrontPort
	RearPorts  map[string]*RearPort

	MetaData *MetaData
}

func NewDevice(name string) *Device {
	return &Device{
		Name:       name,
		Interfaces: make(map[string]*Interface),
		FrontPorts: make(map[string]*FrontPort),
		RearPorts:  make(map[string]*RearPort),
		MetaData:   NewMetaData(),
	}
}

func (d *Device) AddInterfaceItNotExists(ifName string) *Interface {
	iface, exists := d.Interfaces[ifName]
	if !exists {
		iface = newInterface(ifName)
		d.Interfaces[ifName] = iface
	}

	return iface
}

func (d *Device) GetInterface(ifName string) *Interface {
	return d.Interfaces[ifName]
}

func (d *Device) ToProto() *octopuspb.Device {
	if d == nil {
		return nil
	}

	protoDev := &octopuspb.Device{
		Name:       d.Name,
		Status:     d.Status,
		Role:       d.Role,
		Platform:   d.Platform,
		DeviceType: d.DeviceType,

		MetaData: d.MetaData.ToProto(),
	}

	if d.Colo != nil {
		protoDev.ColoId = int32(d.Colo.Id)
	}

	if d.Site != nil {
		protoDev.SiteName = d.Site.Name
	}

	if len(d.Interfaces) > 0 {
		protoDev.Interfaces = make([]*octopuspb.Interface, 0)
		for _, iface := range d.Interfaces {
			protoDev.Interfaces = append(protoDev.Interfaces, iface.ToProto())
		}
	}

	if len(d.FrontPorts) > 0 {
		protoDev.FrontPorts = make([]*octopuspb.FrontPort, 0)
		for _, fp := range d.FrontPorts {
			protoDev.FrontPorts = append(protoDev.FrontPorts, fp.ToProto())
		}
	}

	if len(d.RearPorts) > 0 {
		protoDev.RearPorts = make([]*octopuspb.RearPort, 0)
		for _, rp := range d.RearPorts {
			protoDev.RearPorts = append(protoDev.RearPorts, rp.ToProto())
		}
	}

	return protoDev
}
