package model

const TableNameDcimDevice = "dcim_device"

// DcimDevice mapped from table <dcim_device>
type DcimDevice struct {
	//	Created          time.Time `gorm:"column:created" json:"created"`
	//	LastUpdated      time.Time `gorm:"column:last_updated" json:"last_updated"`
	//  CustomFieldData string `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	//LocalContextData string `gorm:"column:local_context_data" json:"local_context_data"`
	Name     string `gorm:"column:name" json:"name"`
	Serial   string `gorm:"column:serial;not null" json:"serial"`
	AssetTag string `gorm:"column:asset_tag" json:"asset_tag"`
	Position int16  `gorm:"column:position" json:"position"`
	//	Face             string         `gorm:"column:face;not null" json:"face"`
	Status string `gorm:"column:status;not null" json:"status"`
	//	VcPosition       int16          `gorm:"column:vc_position" json:"vc_position"`
	//	VcPriority       int16          `gorm:"column:vc_priority" json:"vc_priority"`
	Comments string `gorm:"column:comments;not null" json:"comments"`
	//	ClusterID        int64          `gorm:"column:cluster_id" json:"cluster_id"`
	DeviceRoleID int64 `gorm:"column:device_role_id;not null" json:"device_role_id"`
	DeviceTypeID int64 `gorm:"column:device_type_id;not null" json:"device_type_id"`
	LocationID   int64 `gorm:"column:location_id" json:"location_id"`
	PlatformID   int64 `gorm:"column:platform_id" json:"platform_id"`
	PrimaryIp4ID int64 `gorm:"column:primary_ip4_id" json:"primary_ip4_id" sql:"primary_ip4_id"`
	PrimaryIp6ID int64 `gorm:"column:primary_ip6_id" json:"primary_ip6_id" sql:"primary_ip6_id"`
	RackID       int64 `gorm:"column:rack_id" json:"rack_id"`
	SiteID       int64 `gorm:"column:site_id;not null" json:"site_id"`
	TenantID     int64 `gorm:"column:tenant_id" json:"tenant_id"`
	//	VirtualChassisID int64          `gorm:"column:virtual_chassis_id" json:"virtual_chassis_id"`
	//	Airflow          string         `gorm:"column:airflow;not null" json:"airflow"`
	DeviceRole DcimDevicerole `pg:"fk:device_role_id"`
	DeviceType DcimDevicetype `pg:"fk:device_type_id"`
	Site       DcimSite       `pg:"fk:site_id"`
	Tags       []string       `sql:"-"`
}

// TableName DcimDevice's table name
func (*DcimDevice) TableName() string {
	return TableNameDcimDevice
}
