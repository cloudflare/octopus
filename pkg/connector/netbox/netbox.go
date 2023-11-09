//
// Copyright (c) 2023 Cloudflare, Inc.
//
// Licensed under Apache 2.0 license found in the LICENSE file
// or at http://www.apache.org/licenses/LICENSE-2.0
//

package netbox

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	dbModel "github.com/cloudflare/octopus/pkg/connector/netbox/model"
	nbUtils "github.com/cloudflare/octopus/pkg/connector/netbox/utils"
	"github.com/cloudflare/octopus/pkg/model"
	"github.com/cloudflare/octopus/proto/octopus"

	bnet "github.com/bio-routing/bio-rd/net"
	log "github.com/sirupsen/logrus"
)

const (
	connectorName  = "Netbox"
	updateInterval = time.Minute * 2
)

type NetboxConnector struct {
	connectorMu       sync.RWMutex
	client            NetboxClientI
	loadDuration      time.Duration
	loadTime          time.Time
	refreshErrorCount atomic.Uint64

	devices             map[int64]*dbModel.DcimDevice
	interfaces          map[int64]*dbModel.DcimInterface
	ipAddresses         []*dbModel.IpamIpaddress
	cables              []*dbModel.DcimCable
	prefixes            []*dbModel.IpamPrefix
	circuits            map[int64]*dbModel.CircuitsCircuit
	circuitTerminations map[int64]*dbModel.CircuitsCircuittermination
	frontPorts          map[int64]*dbModel.DcimFrontport
	rearPorts           map[int64]*dbModel.DcimRearport
}

type NetboxClientI interface {
	Connect() error
	GetDBHost() string
	GetDevices() ([]*dbModel.DcimDevice, error)
	GetInterfaces() (map[int64]*dbModel.DcimInterface, error)
	GetIPAddresses() ([]*dbModel.IpamIpaddress, error)
	GetCables() ([]*dbModel.DcimCable, error)
	GetPrefixes() ([]*dbModel.IpamPrefix, error)
	GetDcimInterfaceTypeID() int32
	GetCircuitsCircuitterminationTypeID() int32
	GetDcimFrontPortTypeID() int32
	GetDcimRearPortTypeID() int32
	GetCircuits() ([]*dbModel.CircuitsCircuit, error)
	GetCircuitTerminations() ([]*dbModel.CircuitsCircuittermination, error)
	GetFrontPorts() ([]*dbModel.DcimFrontport, error)
	GetRearPorts() ([]*dbModel.DcimRearport, error)
}

func NewConnector(host string, port uint, user string, password string, dbName string, useTLS bool, caCertPath string, logDBQueries bool) *NetboxConnector {
	dbParams := dbParams{
		host:         host,
		port:         port,
		dBname:       dbName,
		user:         user,
		password:     password,
		useTLS:       useTLS,
		caCertPath:   caCertPath,
		logDBQueries: logDBQueries,
	}

	return newNetboxConnectorWithClient(&NetboxClient{
		db: newDB(dbParams),
	})
}

func newNetboxConnectorWithClient(apiClient NetboxClientI) *NetboxConnector {
	return &NetboxConnector{
		client: apiClient,
	}
}

func (n *NetboxConnector) InitialLoad() error {
	return n.update()
}

func (n *NetboxConnector) Healthy() bool {
	n.connectorMu.RLock()
	defer n.connectorMu.RUnlock()

	return n._healthy()
}

func (n *NetboxConnector) _healthy() bool {
	return len(n.devices) > 0
}

func (n *NetboxConnector) GetLoadDuration() time.Duration {
	return n.loadDuration
}

func (n *NetboxConnector) GetLoadTime() time.Time {
	return n.loadTime
}

func (n *NetboxConnector) GetName() string {
	return connectorName
}

func (n *NetboxConnector) GetUpdateErrorCount() uint64 {
	return n.refreshErrorCount.Load()
}

func (n *NetboxConnector) EnrichTopology(t *model.Topology) error {
	n.connectorMu.RLock()
	defer n.connectorMu.RUnlock()

	if !n._healthy() {
		return fmt.Errorf("%s not healthy", connectorName)
	}

	return n._enrichTopology(t)
}

func (n *NetboxConnector) _enrichTopology(t *model.Topology) error {
	err := n.addDevices(t)
	if err != nil {
		return fmt.Errorf("failed to enrich devices: %v", err)
	}

	err = n.addInterfaces(t)
	if err != nil {
		return fmt.Errorf("failed to enrich interfaces: %v", err)
	}

	err = n.addInterfaceUnits(t)
	if err != nil {
		return fmt.Errorf("failed to enrich interface units: %v", err)
	}

	err = n.addIPAddresses(t)
	if err != nil {
		return fmt.Errorf("failed to enrich IP addresses: %v", err)
	}

	err = n.addCables(t)
	if err != nil {
		return fmt.Errorf("failed to enrich cables: %v", err)
	}

	err = n.addCircuits(t)
	if err != nil {
		return fmt.Errorf("failed to enrich circuits: %v", err)
	}

	err = n.addRearPorts(t)
	if err != nil {
		return fmt.Errorf("failed to enrich rear ports: %v", err)
	}

	err = n.addFrontPorts(t)
	if err != nil {
		return fmt.Errorf("failed to enrich front ports: %v", err)
	}

	err = n.addPrefixes(t)
	if err != nil {
		return fmt.Errorf("failed to add prefixes: %v", err)
	}

	return nil
}

func (n *NetboxConnector) addDevices(t *model.Topology) error {
	for _, d := range n.devices {
		topoDev := t.AddDeviceIfNotExists(d.Name)
		s := t.AddSiteIfNotExists(d.Site.Name)

		topoDev.Site = s
		topoDev.Role = d.DeviceRole.Slug
		topoDev.DeviceType = d.DeviceType.Slug

		md, err := nbUtils.GetMetaDataFromTags(d.Tags)
		if err != nil {
			return fmt.Errorf("unable to get meta data: %v", err)
		}

		topoDev.MetaData = md
	}

	return nil
}

func (n *NetboxConnector) addInterfaces(t *model.Topology) error {
	for _, nbIfa := range n.interfaces {
		if nbIfa.Parent != nil {
			continue
		}

		d := t.GetDevice(nbIfa.Device.Name)
		if d == nil {
			return fmt.Errorf("can not find device %q", nbIfa.Device.Name)
		}

		ifa := d.AddInterfaceItNotExists(nbIfa.Name)
		t.DevicesByInterfaceID[nbIfa.ID] = d
		t.Interfaces[nbIfa.ID] = ifa

		md, err := nbUtils.GetMetaDataFromTags(nbIfa.Tags)
		if err != nil {
			return fmt.Errorf("unable to get meta data: %v", err)
		}

		ifa.MetaData = md
		ifa.Type = nbIfa.Type

		if nbIfa.LAG != nil {
			ifa.LAGMemberOf = nbIfa.LAG.Name
		}
	}

	return nil
}

func (n *NetboxConnector) addInterfaceUnits(t *model.Topology) error {
	for _, nbIfa := range n.interfaces {
		if nbIfa.Parent == nil {
			continue
		}

		if !strings.Contains(nbIfa.Name, ".") {
			continue
		}

		unitStr := strings.TrimPrefix(nbIfa.Name, nbIfa.Parent.Name+".")
		d := t.GetDevice(nbIfa.Device.Name)
		if d == nil {
			return fmt.Errorf("can not find device %q", nbIfa.Device.Name)
		}

		ifa := d.GetInterface(nbIfa.Parent.Name)
		if ifa == nil {
			return fmt.Errorf("can not find interface %s:%s", nbIfa.Device.Name, nbIfa.Parent.Name)
		}

		vlanTag, err := nbUtils.ParseUnitStr(unitStr)
		if err != nil {
			return fmt.Errorf("unable to convert unit %q (id=%d) (interface %q) to int for %s:%s. Ignoring logical interface", unitStr, nbIfa.ID, nbIfa.Name, nbIfa.Device.Name, nbIfa.Parent.Name)
		}

		u := ifa.AddUnitIfNotExists(vlanTag)
		t.DevicesByInterfaceID[nbIfa.ID] = d

		md, err := nbUtils.GetMetaDataFromTags(nbIfa.Tags)
		if err != nil {
			return fmt.Errorf("unable to get meta data: %v", err)
		}

		u.MetaData = md
	}

	return nil
}

func (n *NetboxConnector) addIPAddresses(t *model.Topology) error {
	for _, nbIP := range n.ipAddresses {
		if nbIP.AssignedObjectID == 0 || nbIP.AssignedObjectTypeID != n.client.GetDcimInterfaceTypeID() {
			continue
		}

		dcimIfa, exists := n.interfaces[nbIP.AssignedObjectID]
		if !exists {
			continue
		}

		ifaID := nbIP.AssignedObjectID
		if dcimIfa.ParentID != 0 {
			ifaID = dcimIfa.ParentID
		}

		ifa := t.Interfaces[ifaID]
		if ifa == nil {
			return fmt.Errorf("interface with id %d not found", ifaID)
		}

		_, vt, err := nbUtils.GetInterfaceAndVLANTag(dcimIfa.Name)
		if err != nil {
			return fmt.Errorf("unable to extract interface name and unit from %q: %v", dcimIfa.Name, err)
		}

		pfx, err := bnet.PrefixFromString(nbUtils.SanitizeIPAddress(nbIP.Address))
		if err != nil {
			return fmt.Errorf("failed to parse IP %q: %v", nbIP.Address, err)
		}

		ip := model.NewIP(*pfx)
		nbUtils.GetCustomFieldData(ip.MetaData, nbIP.CustomFieldData)

		u := ifa.AddUnitIfNotExists(vt)
		if pfx.Addr().IsIPv4() {
			u.IPv4Addresses = append(u.IPv4Addresses, ip)
		} else {
			u.IPv6Addresses = append(u.IPv6Addresses, ip)
		}
	}

	return nil
}

func (n *NetboxConnector) addPrefixes(t *model.Topology) error {
	for _, p := range n.prefixes {
		pfx, err := bnet.PrefixFromString(p.Prefix)
		if err != nil {
			return fmt.Errorf("failed to parse Prefix %q: %v", p.Prefix, err)
		}

		md, err := nbUtils.GetMetaDataFromTags(p.Tags)
		if err != nil {
			return fmt.Errorf("failed to get Tags for Prefix %q: %v", p.Prefix, err)
		}

		oPfx := model.NewPrefix(*pfx)
		oPfx.MetaData = md
		t.Prefixes[p.ID] = oPfx

	}

	return nil
}

func (n *NetboxConnector) getCableEnd(terminationType int32, terminationID int64, t *model.Topology) (*model.CableEnd, error) {
	ce := model.CableEnd{}

	switch terminationType {
	case n.client.GetDcimInterfaceTypeID():
		ce.EndpointType = octopus.CableEndpointType_CABLE_ENDPOINT_TYPE_INTERFACE
		ifa, exists := t.Interfaces[terminationID]
		if !exists {
			return nil, fmt.Errorf("unable to find interface %d", terminationID)
		}

		dev, exists := t.DevicesByInterfaceID[terminationID]
		if !exists {
			return nil, fmt.Errorf("unable to find device by interface id %d", terminationID)
		}

		ce.DeviceName = dev.Name
		ce.EndpointName = ifa.Name

	case n.client.GetCircuitsCircuitterminationTypeID():
		ce.EndpointType = octopus.CableEndpointType_CABLE_ENDPOINT_TYPE_CIRCUIT_TERMINATION

		cktTerm := n.circuitTerminations[terminationID]
		if cktTerm == nil {
			return nil, fmt.Errorf("unable to find circuit termination %d", terminationID)
		}

		ckt := n.circuits[cktTerm.CircuitID]
		if ckt == nil {
			return nil, fmt.Errorf("unable to find circuit %d", cktTerm.ID)
		}

		ce.DeviceName = ckt.Cid
		if ckt.TerminationAID == terminationID {
			ce.EndpointName = "A"
		}

		if ckt.TerminationZID == terminationID {
			ce.EndpointName = "Z"
		}

	case n.client.GetDcimFrontPortTypeID():
		ce.EndpointType = octopus.CableEndpointType_CABLE_ENDPOINT_TYPE_FRONT_PORT
		fp := n.frontPorts[terminationID]
		if fp == nil {
			return nil, fmt.Errorf("unable to find front port %d", terminationID)
		}

		d := n.devices[fp.DeviceID]
		if d == nil {
			return nil, fmt.Errorf("unable to find device with id %d", fp.DeviceID)
		}

		ce.DeviceName = d.Name
		ce.EndpointName = fp.Name

	case n.client.GetDcimRearPortTypeID():
		ce.EndpointType = octopus.CableEndpointType_CABLE_ENDPOINT_TYPE_REAR_PORT
		rp := n.rearPorts[terminationID]
		if rp == nil {
			return nil, fmt.Errorf("unable to find rear port %d", terminationID)
		}

		d := n.devices[rp.DeviceID]
		if d == nil {
			return nil, fmt.Errorf("unable to find device with id %d", rp.DeviceID)
		}

		ce.DeviceName = d.Name
		ce.EndpointName = rp.Name

	default:
		// consoleport, consoleserverport, powerport, or poweroutlet
		return nil, fmt.Errorf("don't know what to do with cable termination ID %d (type %d)", terminationID, terminationType)
	}

	return &ce, nil
}

func (n *NetboxConnector) addCables(t *model.Topology) error {
	for _, c := range n.cables {
		if c.TerminationAID == 0 || c.TerminationBID == 0 {
			continue
		}

		AEnd, err := n.getCableEnd(c.TerminationATypeID, c.TerminationAID, t)
		if err != nil {
			continue
		}

		BEnd, err := n.getCableEnd(c.TerminationBTypeID, c.TerminationBID, t)
		if err != nil {
			continue
		}

		cable := model.Cable{
			AEnd: *AEnd,
			BEnd: *BEnd,
		}

		t.Cables[cable.String()] = &cable
	}

	return nil
}

func (n *NetboxConnector) addCircuits(t *model.Topology) error {
	for _, c := range n.circuits {
		t.Circuits[c.Cid] = model.NewCircuit(c.Cid, c.Provider.Slug, c.Type.Slug, c.Status)
	}

	return nil
}

func (n *NetboxConnector) addRearPorts(t *model.Topology) error {
	for _, rp := range n.rearPorts {
		nbDev := n.devices[rp.DeviceID]
		if nbDev == nil {
			return fmt.Errorf("device %d not found", rp.DeviceID)
		}

		d := t.GetDevice(nbDev.Name)
		if d == nil {
			return fmt.Errorf("can not find device %q", nbDev.Name)
		}

		d.RearPorts[rp.Name] = &model.RearPort{
			Name:      rp.Name,
			Positions: rp.Positions,
		}
	}

	return nil
}

func (n *NetboxConnector) addFrontPorts(t *model.Topology) error {
	for _, fp := range n.frontPorts {
		nbDev := n.devices[fp.DeviceID]
		if nbDev == nil {
			return fmt.Errorf("device %d not found", fp.DeviceID)
		}

		d := t.GetDevice(nbDev.Name)
		if d == nil {
			return fmt.Errorf("can not find device %q", nbDev.Name)
		}

		rpName := ""
		rp := n.rearPorts[fp.RearPortID]
		if rp != nil {
			rpName = rp.Name
		}

		d.FrontPorts[fp.Name] = &model.FrontPort{
			Name:             fp.Name,
			RearPort:         rpName,
			RearPortPosition: uint32(fp.RearPortPosition),
		}
	}

	return nil
}

func (n *NetboxConnector) StartRefreshRoutine() {
	go n.refreshRoutine()
}

func (n *NetboxConnector) refreshRoutine() {
	ticker := time.NewTicker(updateInterval)
	for {
		err := n.update()
		if err != nil {
			n.refreshErrorCount.Add(1)
			log.Errorf("Failed to refresh Netbox data: %v", err)
		} else {
			log.Infof("Successfully refreshed Netbox data from %q", n.client.GetDBHost())
		}

		<-ticker.C
	}
}

func (n *NetboxConnector) update() error {
	startTime := time.Now()
	err := n.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	devices, err := n.client.GetDevices()
	if err != nil {
		return fmt.Errorf("unable to get devices: %v", err)
	}

	interfaces, err := n.client.GetInterfaces()
	if err != nil {
		return fmt.Errorf("unable to get interfaces: %v", err)
	}

	ips, err := n.client.GetIPAddresses()
	if err != nil {
		return fmt.Errorf("unable to get IP addresses: %v", err)
	}

	cables, err := n.client.GetCables()
	if err != nil {
		return fmt.Errorf("unable to get cables: %v", err)
	}

	prefixes, err := n.client.GetPrefixes()
	if err != nil {
		return fmt.Errorf("unable to get prefixes: %v", err)
	}

	circuits, err := n.client.GetCircuits()
	if err != nil {
		return fmt.Errorf("unable to get circuits: %v", err)
	}

	cts, err := n.client.GetCircuitTerminations()
	if err != nil {
		return fmt.Errorf("unable to get circuit terminations: %v", err)
	}

	fps, err := n.client.GetFrontPorts()
	if err != nil {
		return fmt.Errorf("unable to get front ports: %v", err)
	}

	rps, err := n.client.GetRearPorts()
	if err != nil {
		return fmt.Errorf("unable to get rear ports: %v", err)
	}

	n.connectorMu.Lock()
	defer n.connectorMu.Unlock()

	n.devices = make(map[int64]*dbModel.DcimDevice)
	for _, d := range devices {
		n.devices[d.ID] = d
	}

	n.interfaces = interfaces
	n.ipAddresses = ips
	n.cables = cables
	n.prefixes = prefixes
	n.circuits = make(map[int64]*dbModel.CircuitsCircuit)
	for _, ckt := range circuits {
		n.circuits[ckt.ID] = ckt
	}

	n.circuitTerminations = make(map[int64]*dbModel.CircuitsCircuittermination)
	for _, ct := range cts {
		n.circuitTerminations[ct.ID] = ct
	}

	n.frontPorts = make(map[int64]*dbModel.DcimFrontport)
	for _, fp := range fps {
		n.frontPorts[fp.ID] = fp
	}

	n.rearPorts = make(map[int64]*dbModel.DcimRearport)
	for _, rp := range rps {
		n.rearPorts[rp.ID] = rp
	}

	n.loadDuration = time.Since(startTime)
	n.loadTime = time.Now()

	return nil
}
