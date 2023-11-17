package type_configs

import (
	"database/sql/driver"
	"encoding/json"
)

type FilterConfig struct {
	WhereQuery string `json:"whereQuery"`
	Limit      string `json:"limit"`
}

type SourceConfig struct {
	TargetCollectionName string       `json:"targetCollectionName" validate:"required"`
	TargetDatabaseName   string       `json:"targetDatabaseName" validate:"required"`
	DbSchema             string       `json:"dbSchema" validate:"required"`
	TableName            string       `json:"tableName" validate:"required"`
	PrimaryKey           string       `json:"primaryKey" validate:"required"`
	PrimaryKeyType       string       `json:"primaryKeyType" validate:"required"`
	ColumnList           []string     `json:"columnList" validate:"required"`
	FilterConfig         FilterConfig `json:"filterConfig" validate:"required"`
}

type SourceConfigsDto struct {
	Identifier   string         `json:"identifier" validate:"required"`
	SourceConfig []SourceConfig `json:"sourceConfig" gorm:"type:jsonb;default:'[]';not null" validate:"required"`
}

type SourceConfigs struct {
	Identifier   string `json:"identifier" validate:"required" gorm:"primaryKey,index"`
	SourceConfig JSONB  `json:"sourceConfig" gorm:"type:jsonb;default:'[]';not null" validate:"required"`
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}
