// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameDcimInterfaceTaggedVlan = "dcim_interface_tagged_vlans"

// DcimInterfaceTaggedVlan mapped from table <dcim_interface_tagged_vlans>
type DcimInterfaceTaggedVlan struct {
	ID          int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	InterfaceID int64 `gorm:"column:interface_id;not null" json:"interface_id"`
	VlanID      int64 `gorm:"column:vlan_id;not null" json:"vlan_id"`
}

// TableName DcimInterfaceTaggedVlan's table name
func (*DcimInterfaceTaggedVlan) TableName() string {
	return TableNameDcimInterfaceTaggedVlan
}
