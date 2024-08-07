// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCircuitsProvider = "circuits_provider"

// CircuitsProvider mapped from table <circuits_provider>
type CircuitsProvider struct {
	ID              int64    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	// Created         time.Time `gorm:"column:created" json:"created"`
	// LastUpdated     time.Time `gorm:"column:last_updated" json:"last_updated"`
	Name            string    `gorm:"column:name;not null" json:"name"`
	Slug            string    `gorm:"column:slug;not null" json:"slug"`
	// Comments        string    `gorm:"column:comments;not null" json:"comments"`
	// CustomFieldData string    `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	// Description     string    `gorm:"column:description;not null" json:"description"`
}

// TableName CircuitsProvider's table name
func (*CircuitsProvider) TableName() string {
	return TableNameCircuitsProvider
}
