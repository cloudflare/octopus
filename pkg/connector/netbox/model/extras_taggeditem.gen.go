package model

const TableNameExtrasTaggeditem = "extras_taggeditem"

// ExtrasTaggeditem mapped from table <extras_taggeditem>
type ExtrasTaggeditem struct {
	ObjectID      int32     `gorm:"column:object_id;not null" json:"object_id"`
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	ContentTypeID int32     `gorm:"column:content_type_id;not null" json:"content_type_id"`
	TagID         int64     `gorm:"column:tag_id;not null" json:"tag_id"`
	Tag           ExtrasTag `pg:"fk:tag_id"`
}

// TableName ExtrasTaggeditem's table name
func (*ExtrasTaggeditem) TableName() string {
	return TableNameExtrasTaggeditem
}
