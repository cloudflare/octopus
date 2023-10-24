package model

const TableNameDcimRearport = "dcim_rearport"

// DcimRearport mapped from table <dcim_rearport>
type DcimRearport struct {
	//Created         time.Time `gorm:"column:created" json:"created"`
	//LastUpdated     time.Time `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null" json:"name"`
	//Label           string    `gorm:"column:label;not null" json:"label"`
	//Description     string    `gorm:"column:description;not null" json:"description"`
	//LinkPeerID      int64     `gorm:"column:_link_peer_id" json:"_link_peer_id"`
	//MarkConnected   bool      `gorm:"column:mark_connected;not null" json:"mark_connected"`
	//Type            string    `gorm:"column:type;not null" json:"type"`
	Positions int16 `gorm:"column:positions;not null" json:"positions"`
	//LinkPeerTypeID  int32     `gorm:"column:_link_peer_type_id" json:"_link_peer_type_id"`
	//CableID         int64     `gorm:"column:cable_id" json:"cable_id"`
	DeviceID int64 `gorm:"column:device_id;not null" json:"device_id"`
	//Color           string    `gorm:"column:color;not null" json:"color"`
	//ModuleID        int64     `gorm:"column:module_id" json:"module_id"`
}

// TableName DcimRearport's table name
func (*DcimRearport) TableName() string {
	return TableNameDcimRearport
}
