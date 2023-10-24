package model

const TableNameIpamIpaddress = "ipam_ipaddress"

// IpamIpaddress mapped from table <ipam_ipaddress>
type IpamIpaddress struct {
	//Created              time.Time `gorm:"column:created" json:"created"`
	//LastUpdated          time.Time `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData string `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Address string `gorm:"column:address;not null" json:"address"`
	//Status               string `gorm:"column:status;not null" json:"status"`
	//Role                 string `gorm:"column:role;not null" json:"role"`
	AssignedObjectID int64 `gorm:"column:assigned_object_id" json:"assigned_object_id"`
	//DNSName              string `gorm:"column:dns_name;not null" json:"dns_name"`
	//Description          string `gorm:"column:description;not null" json:"description"`
	AssignedObjectTypeID int32 `gorm:"column:assigned_object_type_id" json:"assigned_object_type_id"`
	//NatInsideID          int64     `gorm:"column:nat_inside_id" json:"nat_inside_id"`
	//TenantID int64 `gorm:"column:tenant_id" json:"tenant_id"`
	//VrfID    int64 `gorm:"column:vrf_id" json:"vrf_id"`
}

// TableName IpamIpaddress's table name
func (*IpamIpaddress) TableName() string {
	return TableNameIpamIpaddress
}
