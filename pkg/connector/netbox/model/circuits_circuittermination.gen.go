package model

const TableNameCircuitsCircuittermination = "circuits_circuittermination"

// CircuitsCircuittermination mapped from table <circuits_circuittermination>
type CircuitsCircuittermination struct {
	//Created     time.Time `gorm:"column:created" json:"created"`
	//LastUpdated time.Time `gorm:"column:last_updated" json:"last_updated"`
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	//LinkPeerID        int64     `gorm:"column:_link_peer_id" json:"_link_peer_id"`
	//MarkConnected bool   `gorm:"column:mark_connected;not null" json:"mark_connected"`
	//TermSide      string `gorm:"column:term_side;not null" json:"term_side"`
	//PortSpeed     int32  `gorm:"column:port_speed" json:"port_speed"`
	//UpstreamSpeed int32  `gorm:"column:upstream_speed" json:"upstream_speed"`
	//XconnectID    string `gorm:"column:xconnect_id;not null" json:"xconnect_id"`
	//PpInfo        string `gorm:"column:pp_info;not null" json:"pp_info"`
	//Description   string `gorm:"column:description;not null" json:"description"`
	//LinkPeerTypeID    int32     `gorm:"column:_link_peer_type_id" json:"_link_peer_type_id"`
	//CableID           int64 `gorm:"column:cable_id" json:"cable_id"`
	CircuitID int64 `gorm:"column:circuit_id;not null" json:"circuit_id"`
	//ProviderNetworkID int64 `gorm:"column:provider_network_id" json:"provider_network_id"`
	//SiteID            int64 `gorm:"column:site_id" json:"site_id"`
}

// TableName CircuitsCircuittermination's table name
func (*CircuitsCircuittermination) TableName() string {
	return TableNameCircuitsCircuittermination
}
