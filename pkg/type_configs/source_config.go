package type_configs

type FilterConfig struct {
	WhereQuery string `json:"whereQuery"`
	Limit      string `json:"limit"`
}

type SourceConfig struct {
	TargetCollectionName string       `json:"targetCollectionName" validate:"required"`
	DbSchema             string       `json:"dbSchema" validate:"required"`
	TableName            string       `json:"tableName" validate:"required"`
	PrimaryKey           string       `json:"primaryKey" validate:"required"`
	PrimaryKeyType       string       `json:"primaryKeyType" validate:"required"`
	ColumnList           []string     `json:"columnList" validate:"required"`
	FilterConfig         FilterConfig `json:"filterConfig" validate:"required"`
}

type PollerConfig struct {
	PollingStrategy string `json:"pollingStrategy"`
}

type SourceConfigsDto struct {
	Identifier   string         `json:"identifier" validate:"required"`
	SourceConfig []SourceConfig `json:"sourceConfig" gorm:"type:jsonb;default:'[]';not null" validate:"required"`
	PollerConfig PollerConfig   `json:"pollerConfig" validate:"required"`
}

type SourceConfigs struct {
	Identifier   string         `json:"identifier" validate:"required" gorm:"primaryKey,index"`
	SourceConfig []SourceConfig `json:"sourceConfig" gorm:"serializer:json;not null" validate:"required"`
	Type         string         `json:"type"`
	PollerConfig PollerConfig   `json:"pollerConfig" gorm:"serializer:json" validate:"required"`
}
