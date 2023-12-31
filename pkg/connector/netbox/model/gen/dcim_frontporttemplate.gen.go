// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameDcimFrontporttemplate = "dcim_frontporttemplate"

// DcimFrontporttemplate mapped from table <dcim_frontporttemplate>
type DcimFrontporttemplate struct {
	Created          time.Time `gorm:"column:created" json:"created"`
	LastUpdated      time.Time `gorm:"column:last_updated" json:"last_updated"`
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name             string    `gorm:"column:name;not null" json:"name"`
	Label            string    `gorm:"column:label;not null" json:"label"`
	Description      string    `gorm:"column:description;not null" json:"description"`
	Type             string    `gorm:"column:type;not null" json:"type"`
	RearPortPosition int16     `gorm:"column:rear_port_position;not null" json:"rear_port_position"`
	DeviceTypeID     int64     `gorm:"column:device_type_id" json:"device_type_id"`
	RearPortID       int64     `gorm:"column:rear_port_id;not null" json:"rear_port_id"`
	Color            string    `gorm:"column:color;not null" json:"color"`
	ModuleTypeID     int64     `gorm:"column:module_type_id" json:"module_type_id"`
}

// TableName DcimFrontporttemplate's table name
func (*DcimFrontporttemplate) TableName() string {
	return TableNameDcimFrontporttemplate
}
