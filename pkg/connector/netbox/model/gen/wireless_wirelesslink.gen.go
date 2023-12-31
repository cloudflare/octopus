// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameWirelessWirelesslink = "wireless_wirelesslink"

// WirelessWirelesslink mapped from table <wireless_wirelesslink>
type WirelessWirelesslink struct {
	Created            time.Time `gorm:"column:created" json:"created"`
	LastUpdated        time.Time `gorm:"column:last_updated" json:"last_updated"`
	CustomFieldData    string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID                 int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Ssid               string    `gorm:"column:ssid;not null" json:"ssid"`
	Status             string    `gorm:"column:status;not null" json:"status"`
	Description        string    `gorm:"column:description;not null" json:"description"`
	AuthCipher         string    `gorm:"column:auth_cipher;not null" json:"auth_cipher"`
	AuthPsk            string    `gorm:"column:auth_psk;not null" json:"auth_psk"`
	AuthType           string    `gorm:"column:auth_type;not null" json:"auth_type"`
	InterfaceADeviceID int64     `gorm:"column:_interface_a_device_id" json:"_interface_a_device_id"`
	InterfaceBDeviceID int64     `gorm:"column:_interface_b_device_id" json:"_interface_b_device_id"`
	InterfaceAID       int64     `gorm:"column:interface_a_id;not null" json:"interface_a_id"`
	InterfaceBID       int64     `gorm:"column:interface_b_id;not null" json:"interface_b_id"`
}

// TableName WirelessWirelesslink's table name
func (*WirelessWirelesslink) TableName() string {
	return TableNameWirelessWirelesslink
}
