package pocketsql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/evantbyrne/trance"
	"github.com/evantbyrne/trance/sqlitedialect"

	// Future:
	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"
)

func connectionInfo(connection string) (string, string) {
	halves := strings.SplitN(connection, ":", 2)
	if strings.HasPrefix(halves[1], "//") {
		// Handle connection strings that include driver name as prefix.
		// e.g., postgres://user:password@localhost/db?sslmode=verify-full
		return halves[0], fmt.Sprint(halves[0], ":", halves[1])
	}
	return halves[0], halves[1]
}

func database(connection string, callback func() error) error {
	driverName, dataSourceName := connectionInfo(connection)
	switch driverName {
	// Future:
	// case "mysql":
	// 	trance.SetDialect(mysqldialect.MysqlDialect{})
	// case "postgres":
	// 	trance.SetDialect(pqdialect.PqDialect{})
	case "sqlite3":
		trance.SetDialect(sqlitedialect.SqliteDialect{})
	default:
		return trance.ErrorInternalServer{Message: fmt.Sprintf("Unsupported driver '%s' provided.", driverName)}
	}
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	trance.UseDatabase(db)

	return callback()
}

func scanRows(fields map[string]reflect.StructField, headers []string, res *sql.Rows) (int, []map[string]any, []map[string]any, error) {
	n := 0
	modals := make([]map[string]any, 0)
	rows := make([]map[string]any, 0)
	for res.Next() {
		row, err := trance.ScanFieldsToMap(res, fields)
		if err != nil {
			return n, modals, rows, err
		}
		row["modal_id"] = fmt.Sprint("row-", n)
		rows = append(rows, row)

		details := make([]map[string]any, 0)
		for _, name := range headers {
			details = append(details, map[string]any{
				"name":  name,
				"value": row[name],
			})
		}

		modals = append(modals, map[string]any{
			"id":      row["modal_id"],
			"table":   "Custom",
			"details": details,
		})
		n++
	}
	return n, modals, rows, nil
}
