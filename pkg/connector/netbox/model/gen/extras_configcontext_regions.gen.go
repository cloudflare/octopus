// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameExtrasConfigcontextRegion = "extras_configcontext_regions"

// ExtrasConfigcontextRegion mapped from table <extras_configcontext_regions>
type ExtrasConfigcontextRegion struct {
	ID              int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ConfigcontextID int64 `gorm:"column:configcontext_id;not null" json:"configcontext_id"`
	RegionID        int64 `gorm:"column:region_id;not null" json:"region_id"`
}

// TableName ExtrasConfigcontextRegion's table name
func (*ExtrasConfigcontextRegion) TableName() string {
	return TableNameExtrasConfigcontextRegion
}
