// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameExtrasConfigcontextSite = "extras_configcontext_sites"

// ExtrasConfigcontextSite mapped from table <extras_configcontext_sites>
type ExtrasConfigcontextSite struct {
	ID              int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ConfigcontextID int64 `gorm:"column:configcontext_id;not null" json:"configcontext_id"`
	SiteID          int64 `gorm:"column:site_id;not null" json:"site_id"`
}

// TableName ExtrasConfigcontextSite's table name
func (*ExtrasConfigcontextSite) TableName() string {
	return TableNameExtrasConfigcontextSite
}
