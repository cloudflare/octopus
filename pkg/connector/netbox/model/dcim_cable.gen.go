package model

const TableNameDcimCable = "dcim_cable"

// DcimCable mapped from table <dcim_cable>
type DcimCable struct {
	//Created              time.Time `gorm:"column:created" json:"created"`
	//LastUpdated          time.Time `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData      string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID             int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	TerminationAID int64  `gorm:"column:termination_a_id;not null" json:"termination_a_id" sql:"termination_a_id"`
	TerminationBID int64  `gorm:"column:termination_b_id;not null" json:"termination_b_id" sql:"termination_b_id"`
	Type           string `gorm:"column:type;not null" json:"type"`
	Status         string `gorm:"column:status;not null" json:"status"`
	//Label                string    `gorm:"column:label;not null" json:"label"`
	//Color                string    `gorm:"column:color;not null" json:"color"`
	//Length               float64   `gorm:"column:length" json:"length"`
	//LengthUnit           string    `gorm:"column:length_unit;not null" json:"length_unit"`
	//AbsLength            float64   `gorm:"column:_abs_length" json:"_abs_length" sql:"_abs_length"`
	TerminationADeviceID int64 `gorm:"column:_termination_a_device_id" json:"_termination_a_device_id" sql:"_termination_a_device_id"`
	TerminationBDeviceID int64 `gorm:"column:_termination_b_device_id" json:"_termination_b_device_id" sql:"_termination_b_device_id"`
	TerminationATypeID   int32 `gorm:"column:termination_a_type_id;not null" json:"termination_a_type_id" sql:"termination_a_type_id"`
	TerminationBTypeID   int32 `gorm:"column:termination_b_type_id;not null" json:"termination_b_type_id" sql:"termination_b_type_id"`
	TenantID             int64 `gorm:"column:tenant_id" json:"tenant_id"`
}

// TableName DcimCable's table name
func (*DcimCable) TableName() string {
	return TableNameDcimCable
}
