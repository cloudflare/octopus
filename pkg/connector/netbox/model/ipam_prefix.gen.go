package model

const TableNameIpamPrefix = "ipam_prefix"

// IpamPrefix mapped from table <ipam_prefix>
type IpamPrefix struct {
	//Created         time.Time `gorm:"column:created" json:"created"`
	//LastUpdated     time.Time `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Prefix string `gorm:"column:prefix;not null" json:"prefix"`
	//Status string `gorm:"column:status;not null" json:"status"`
	//IsPool          bool      `gorm:"column:is_pool;not null" json:"is_pool"`
	//Description     string    `gorm:"column:description;not null" json:"description"`
	//RoleID int64 `gorm:"column:role_id" json:"role_id"`
	//SiteID          int64     `gorm:"column:site_id" json:"site_id"`
	//TenantID        int64     `gorm:"column:tenant_id" json:"tenant_id"`
	//VlanID          int64     `gorm:"column:vlan_id" json:"vlan_id"`
	//VrfID           int64     `gorm:"column:vrf_id" json:"vrf_id"`
	//MarkUtilized    bool      `gorm:"column:mark_utilized;not null" json:"mark_utilized"`
	Tags []string `sql:"-"`
}

// TableName IpamPrefix's table name
func (*IpamPrefix) TableName() string {
	return TableNameIpamPrefix
}
