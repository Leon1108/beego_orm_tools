package dbschema

import (
	"testing"
)

func TestGetTables(t *testing.T) {
	dsName := "root@/information_schema"
	schema := "test"

	tables := GetAllTableInfos(dsName, schema)
	t.Log("Tables: ", tables)
}
