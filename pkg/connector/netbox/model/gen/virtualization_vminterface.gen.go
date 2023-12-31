// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVirtualizationVminterface = "virtualization_vminterface"

// VirtualizationVminterface mapped from table <virtualization_vminterface>
type VirtualizationVminterface struct {
	Created          time.Time `gorm:"column:created" json:"created"`
	LastUpdated      time.Time `gorm:"column:last_updated" json:"last_updated"`
	CustomFieldData  string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Enabled          bool      `gorm:"column:enabled;not null" json:"enabled"`
	MacAddress       string    `gorm:"column:mac_address" json:"mac_address"`
	Mtu              int32     `gorm:"column:mtu" json:"mtu"`
	Mode             string    `gorm:"column:mode;not null" json:"mode"`
	Name             string    `gorm:"column:name;not null" json:"name"`
	Description      string    `gorm:"column:description;not null" json:"description"`
	ParentID         int64     `gorm:"column:parent_id" json:"parent_id"`
	UntaggedVlanID   int64     `gorm:"column:untagged_vlan_id" json:"untagged_vlan_id"`
	VirtualMachineID int64     `gorm:"column:virtual_machine_id;not null" json:"virtual_machine_id"`
	BridgeID         int64     `gorm:"column:bridge_id" json:"bridge_id"`
	VrfID            int64     `gorm:"column:vrf_id" json:"vrf_id"`
}

// TableName VirtualizationVminterface's table name
func (*VirtualizationVminterface) TableName() string {
	return TableNameVirtualizationVminterface
}
