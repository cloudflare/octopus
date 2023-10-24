package model

const TableNameCircuitsCircuit = "circuits_circuit"

// CircuitsCircuit mapped from table <circuits_circuit>
type CircuitsCircuit struct {
	//Created         time.Time        `gorm:"column:created" json:"created"`
	//LastUpdated     time.Time        `gorm:"column:last_updated" json:"last_updated"`
	//CustomFieldData string           `gorm:"column:custom_field_data;not null" json:"custom_field_data"`
	ID     int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Cid    string `gorm:"column:cid;not null" json:"cid"`
	Status string `gorm:"column:status;not null" json:"status"`
	//InstallDate     time.Time        `gorm:"column:install_date" json:"install_date"`
	//CommitRate      int32            `gorm:"column:commit_rate" json:"commit_rate"`
	//Description     string           `gorm:"column:description;not null" json:"description"`
	//Comments        string           `gorm:"column:comments;not null" json:"comments"`
	ProviderID int64 `gorm:"column:provider_id;not null" json:"provider_id"`
	//TenantID        int64            `gorm:"column:tenant_id" json:"tenant_id"`
	TerminationAID int64               `gorm:"column:termination_a_id" json:"termination_a_id" sql:"termination_a_id"`
	TerminationZID int64               `gorm:"column:termination_z_id" json:"termination_z_id" sql:"termination_z_id"`
	TypeID         int64               `gorm:"column:type_id;not null" json:"type_id"`
	Tags           []string            `sql:"-"`
	Provider       CircuitsProvider    `pg:"fk:provider_id"`
	Type           CircuitsCircuittype `pg:"fk:type_id"`
}

// TableName CircuitsCircuit's table name
func (*CircuitsCircuit) TableName() string {
	return TableNameCircuitsCircuit
}
