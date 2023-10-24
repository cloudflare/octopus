//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package netbox

import (
	"fmt"

	"github.com/cloudflare/octopus/pkg/connector/netbox/model"
)

type NetboxClient struct {
	db *database
}

func (nbc *NetboxClient) Connect() error {
	err := nbc.db.connect()
	if err != nil {
		return err
	}

	return nbc.db.loadContentTypes()
}

func (nbc *NetboxClient) GetDBHost() string {
	return nbc.db.params.host
}

func (nbc *NetboxClient) GetDevices() ([]*model.DcimDevice, error) {
	devices, err := nbc.db.getDevices()
	if err != nil {
		return nil, fmt.Errorf("unable to get devices: %v", err)
	}

	return devices, nil
}

func (nbc *NetboxClient) GetInterfaces() (map[int64]*model.DcimInterface, error) {
	interfaces, err := nbc.db.getInterfaces()
	if err != nil {
		return nil, fmt.Errorf("unable to get interfaces: %v", err)
	}

	return interfaces, nil
}

func (nbc *NetboxClient) GetIPAddresses() ([]*model.IpamIpaddress, error) {
	addrs, err := nbc.db.getIPAddresses()
	if err != nil {
		return nil, fmt.Errorf("unable to get ip addresses: %v", err)
	}

	return addrs, nil
}

func (nbc *NetboxClient) GetCables() ([]*model.DcimCable, error) {
	cables, err := nbc.db.getCables()
	if err != nil {
		return nil, fmt.Errorf("unable to get cables: %v", err)
	}

	return cables, nil
}

func (nbc *NetboxClient) GetPrefixes() ([]*model.IpamPrefix, error) {
	prefixes, err := nbc.db.getPrefixes()
	if err != nil {
		return nil, fmt.Errorf("unable to get prefixes: %v", err)
	}

	return prefixes, nil
}

func (nbc *NetboxClient) GetCircuits() ([]*model.CircuitsCircuit, error) {
	circuits, err := nbc.db.getCircuits()
	if err != nil {
		return nil, fmt.Errorf("unable to get circuits: %v", err)
	}

	return circuits, nil
}

func (nbc *NetboxClient) GetCircuitTerminations() ([]*model.CircuitsCircuittermination, error) {
	cts, err := nbc.db.getCircuitTerminations()
	if err != nil {
		return nil, fmt.Errorf("unable to get circuit terminations: %v", err)
	}

	return cts, nil
}

func (nbc *NetboxClient) GetFrontPorts() ([]*model.DcimFrontport, error) {
	fps, err := nbc.db.getFrontports()
	if err != nil {
		return nil, fmt.Errorf("unable to get front ports: %v", err)
	}

	return fps, nil
}

func (nbc *NetboxClient) GetRearPorts() ([]*model.DcimRearport, error) {
	rps, err := nbc.db.getRearports()
	if err != nil {
		return nil, fmt.Errorf("unable to get front ports: %v", err)
	}

	return rps, nil
}

func (nbc *NetboxClient) GetDcimInterfaceTypeID() int32 {
	return nbc.db.contentTypeDcimInterface
}

func (nbc *NetboxClient) GetCircuitsCircuitterminationTypeID() int32 {
	return nbc.db.contentTypeCircuitsCircuittermination
}

func (nbc *NetboxClient) GetDcimFrontPortTypeID() int32 {
	return nbc.db.contentTypeFrontPort
}

func (nbc *NetboxClient) GetDcimRearPortTypeID() int32 {
	return nbc.db.contentTypeRearPort
}
