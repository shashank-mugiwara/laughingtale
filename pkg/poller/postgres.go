package poller

import (
	"database/sql"
	"encoding/json"

	"github.com/shashank-mugiwara/laughingtale/client"
	"github.com/shashank-mugiwara/laughingtale/logger"
	"github.com/shashank-mugiwara/laughingtale/pkg/type_config"
)

func PollData(sourceConfig *type_config.SourceConfigContainer) {
	var lstOfSources []string
	for _, sc := range sourceConfig.SourceConfig {
		lstOfSources = append(lstOfSources, sc.DbSchema+"."+sc.TableName)
	}

	logger.GetLaughingTaleLogger().Info("List of sources: ", lstOfSources)
	for _, sc := range sourceConfig.SourceConfig {
		go pollDataFromTable(&sc)
	}
}

func pollDataFromTable(sourceConfig *type_config.SourceConfig) {
	listOfFields, err := generateListOfFields(sourceConfig)
	if err != nil {
		logger.GetLaughingTaleLogger().Error("Failed to poll data from database.")
		return
	}

	query := "select " + listOfFields + " FROM " + sourceConfig.DbSchema + "." + sourceConfig.TableName + " WHERE " +
		sourceConfig.FilterConfig.WhereQuery + " LIMIT " + sourceConfig.FilterConfig.Limit + ";"

	logger.GetLaughingTaleLogger().Info("Executing query: ", query)
	rows, err := client.GetPostgresDb().Raw(query).Rows()
	defer rows.Close()
	parse(rows)
}

func generateListOfFields(sourceConfig *type_config.SourceConfig) (string, error) {

	if len(sourceConfig.ColumnList) < 1 {
		logger.GetLaughingTaleLogger().Info("No columns specified, selecting all columns by default")
		return "*", nil
	}

	var selectedFields string = ""
	for _, field := range sourceConfig.ColumnList {
		selectedFields = selectedFields + "," + field
	}

	if len([]rune(selectedFields)) < 2 {
		logger.GetLaughingTaleLogger().Error("Failed to construct query string, so selecting all columns by default")
		return "*", nil
	}

	return selectedFields[1:], nil
}

func parse(rows *sql.Rows) error {
	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {

		scanArgs := make([]interface{}, count)

		for i, v := range columnTypes {

			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(scanArgs...)

		if err != nil {
			return err
		}

		masterData := map[string]interface{}{}

		for i, v := range columnTypes {

			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				masterData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				masterData[v.Name()] = z.String
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				masterData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				masterData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				masterData[v.Name()] = z.Int32
				continue
			}

			masterData[v.Name()] = scanArgs[i]
		}

		finalRows = append(finalRows, masterData)
	}

	z, err := json.Marshal(finalRows)
	logger.GetLaughingTaleLogger().Info(string(z))
	return err
}
