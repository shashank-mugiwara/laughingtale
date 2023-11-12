package dbutils

import "github.com/shashank-mugiwara/laughingtale/client"

func ListTables(schema string) []string {
	var tables []string
	if err := client.GetPostgresDb().Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
		panic(err)
	}
	return tables
}
