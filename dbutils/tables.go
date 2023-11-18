package dbutils

import (
	"github.com/shashank-mugiwara/laughingtale/db"
)

func ListTables(schema string) []string {
	var tables []string
	if err := db.GetlaughingtaleDb().Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		panic(err)
	}
	return tables
}
