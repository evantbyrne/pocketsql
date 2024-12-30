package pocketsql

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/evantbyrne/trance"
)

func customSql(strand *trance.Strand) error {
	q := strings.TrimSpace(strand.Request().FormValue("sql"))
	if q == "" {
		return errorSql{Message: "SQL query required."}
	}
	if strings.ToLower(q[0:6]) == "select" {
		return customSqlQuery(strand, q)
	}
	return customSqlExec(strand, q)
}

func customSqlExec(strand *trance.Strand, q string) error {
	var (
		lastInsertId    int64
		lastInsertIdErr error
		rowsAffected    int64
		rowsAffectedErr error
	)

	res, err := trance.Database().Exec(q)
	if err != nil {
		return errorSql{Message: err.Error()}
	}
	lastInsertId, lastInsertIdErr = res.LastInsertId()
	rowsAffected, rowsAffectedErr = res.RowsAffected()

	links, err := describeSchemaSqlite()
	if err != nil {
		return err
	}

	strand.Response.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	return execStreamTemplate.Execute(strand.Response, map[string]any{
		"lastInsertId":      lastInsertId,
		"lastInsertIdValid": lastInsertIdErr == nil,
		"links":             links,
		"rowsAffected":      rowsAffected,
		"rowsAffectedValid": rowsAffectedErr == nil,
		"sql":               q,
		"table":             "Custom Query",
	})
}

func customSqlQuery(strand *trance.Strand, q string) error {
	res, err := trance.Database().Query(q)
	if err != nil {
		return errorSql{Message: err.Error()}
	}

	columns, err := res.Columns()
	if err != nil {
		return err
	}

	fields := make(map[string]reflect.StructField, 0)
	headers := make([]string, 0)
	for _, column := range columns {
		fields[column] = reflect.StructField{
			Name: column,
			Type: reflect.TypeOf(sql.NullString{}),
		}
		headers = append(headers, column)
	}

	n, modals, rows, err := scanRows(fields, headers, res)
	if err != nil {
		return err
	}

	links, err := describeSchemaSqlite()
	if err != nil {
		return err
	}

	strand.Response.Header().Add("Content-Type", "text/vnd.turbo-stream.html")

	return tableStreamTemplate.Execute(strand.Response, map[string]any{
		"count":   n,
		"end":     n,
		"headers": headers,
		"links":   links,
		"modals":  modals,
		"rows":    rows,
		"sql":     q,
		"start":   1,
		"table":   "Custom Query",
	})
}
