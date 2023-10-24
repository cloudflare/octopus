package model

const TableNameDcimInterface = "dcim_interface"

// DcimInterface mapped from table <dcim_interface>
type DcimInterface struct {
	//Created         time.Time `gorm:"column:created" json:"created"`
	//LastUpdated     time.Time `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData string `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
	//Label       string `gorm:"column:label;not null" json:"label"`
	//Description string `gorm:"column:description;not null" json:"description"`
	//LinkPeerID         int64          `gorm:"column:_link_peer_id" json:"_link_peer_id"`
	//MarkConnected bool   `gorm:"column:mark_connected;not null" json:"mark_connected"`
	//Enabled       bool   `gorm:"column:enabled;not null" json:"enabled"`
	MacAddress string `gorm:"column:mac_address" json:"mac_address"`
	//Mtu           int32  `gorm:"column:mtu" json:"mtu"`
	//Mode          string `gorm:"column:mode;not null" json:"mode"`
	Type     string `gorm:"column:type;not null" json:"type"`
	MgmtOnly bool   `gorm:"column:mgmt_only;not null" json:"mgmt_only"`
	//LinkPeerTypeID     int32          `gorm:"column:_link_peer_type_id" json:"_link_peer_type_id"`
	//PathID             int64          `gorm:"column:_path_id" json:"_path_id"`
	CableID  int64 `gorm:"column:cable_id" json:"cable_id"`
	DeviceID int64 `gorm:"column:device_id;not null" json:"device_id"`
	LagID    int64 `gorm:"column:lag_id" json:"lag_id"`
	ParentID int64 `gorm:"column:parent_id" json:"parent_id"`
	//UntaggedVlanID int64 `gorm:"column:untagged_vlan_id" json:"untagged_vlan_id"`
	//Wwn                string         `gorm:"column:wwn" json:"wwn"`
	//BridgeID           int64          `gorm:"column:bridge_id" json:"bridge_id"`
	//RfRole             string         `gorm:"column:rf_role;not null" json:"rf_role"`
	//RfChannel          string         `gorm:"column:rf_channel;not null" json:"rf_channel"`
	//RfChannelFrequency float64        `gorm:"column:rf_channel_frequency" json:"rf_channel_frequency"`
	//RfChannelWidth     float64        `gorm:"column:rf_channel_width" json:"rf_channel_width"`
	//TxPower            int16          `gorm:"column:tx_power" json:"tx_power"`
	//WirelessLinkID     int64          `gorm:"column:wireless_link_id" json:"wireless_link_id"`
	//ModuleID int64 `gorm:"column:module_id" json:"module_id"`
	//VrfID    int64          `gorm:"column:vrf_id" json:"vrf_id"`
	//Duplex   string         `gorm:"column:duplex" json:"duplex"`
	Speed  int32          `gorm:"column:speed" json:"speed"`
	Parent *DcimInterface `pg:"fk:parent_id"`
	Device DcimDevice     `pg:"fk:device_id"`
	LAG    *DcimInterface `pg:"fk:lag_id"`
	Tags   []string       `sql:"-"`
}

// TableName DcimInterface's table name
func (*DcimInterface) TableName() string {
	return TableNameDcimInterface
}
