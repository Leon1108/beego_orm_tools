package dbschema

import (
	"fmt"
	"strings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type ColumnInfo struct {
	Name	string
	Type	string
	Length	sql.NullInt64
}

type TableInfo struct {
	Name 		string
	Columns		[]ColumnInfo
}

const SCHEMA_SQL = `
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    DATA_TYPE,
    CHARACTER_MAXIMUM_LENGTH
FROM
    INFORMATION_SCHEMA.COLUMNS
WHERE
    TABLE_SCHEMA = ?
ORDER BY TABLE_NAME , ORDINAL_POSITION;
`

func GetAllTableInfos(dsName, schema string) (result []TableInfo){
	db, err := sql.Open("mysql", dsName)
	if nil != err {
		errLog("Fail to open database. Err: %v", err)
		panic(err)
	}
	defer db.Close()

	rows, rowErr := db.Query(SCHEMA_SQL, schema)
	if nil != rowErr {
		errLog("Fail to query. Err: %v", rowErr)
		panic(err)
	}


	var tabs map[string]*TableInfo = make(map[string]*TableInfo)
	for rows.Next() {
		var tableName, columnName, dataType string
		var maxLength sql.NullInt64
		rows.Scan(&tableName, &columnName, &dataType, &maxLength)
		log("%v, %v, %v, %v", tableName, columnName, dataType, maxLength)

		tbName := strings.ToLower(tableName)
		if tbInfo, ok := tabs[tbName]; ok {
			tbInfo.Columns = append(tbInfo.Columns, createColumnInfo(columnName, dataType, maxLength))
		} else { // 还没有这张表
			tabs[tbName] = &TableInfo {
				Name : tableName,
				Columns : []ColumnInfo{createColumnInfo(columnName, dataType, maxLength)},
			}
		}
	}

	for _, v := range tabs {
		result = append(result, *v)
	}

	return
}

func createColumnInfo(columnName, dataType string, maxLength sql.NullInt64) ColumnInfo {
	col := ColumnInfo{
		Name : columnName,
		Type : dataType,
		Length : maxLength,
	}
	return col
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func errLog(format string, args ... interface {}){
	fmt.Printf(format + "\n", args...)
}

func log(format string, args ... interface {}){
	fmt.Printf(format + "\n", args ...)
}
