package type_configs

type FilterConfig struct {
	WhereQuery string `json:"whereQuery"`
	Limit      string `json:"limit"`
}

type SourceConfig struct {
	TargetCollectionName string       `json:"targetCollectionName"`
	TargetDatabaseName   string       `json:"targetDatabaseName"`
	DbSchema             string       `json:"dbSchema"`
	TableName            string       `json:"tableName"`
	PrimaryKey           string       `json:"primaryKey"`
	PrimaryKeyType       string       `json:"primaryKeyType"`
	ColumnList           []string     `json:"columnList"`
	FilterConfig         FilterConfig `json:"filterConfig"`
}

type SourceConfigs struct {
	Identifier   string         `json:"identifier" validate:"required"`
	SourceConfig []SourceConfig `json:"sourceConfig" gorm:"type:jsonb;default:'[]';not null" validate:"required"`
}
