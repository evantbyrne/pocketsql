package pocketsql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"

	"github.com/evantbyrne/trance"
	"github.com/evantbyrne/trance/sqlitedialect"
)

type tableDescriber = func(string) (map[string]reflect.StructField, []string, string, error)

const describeTableLimit = 10

func describeTable(strand *trance.Strand) error {
	var describer tableDescriber
	switch trance.GetDialect() {
	case sqlitedialect.SqliteDialect{}:
		describer = describeTableSqlite
	default:
		return trance.ErrorInternalServer{Message: fmt.Sprintf("Unsupported dialect '%T' for describing tables.", trance.GetDialect())}
	}

	table := strand.Request().PathValue("name")
	fields, headers, pk, err := describer(table)
	if err != nil {
		return err
	}

	offset, err := strconv.Atoi(strand.Request().FormValue("offset"))
	if err != nil {
		offset = 0
	}

	order := strand.Request().FormValue("order")
	if order == "" {
		if pk != "" {
			order = pk
		} else {
			order = headers[0]
		}
	}

	direction := "asc"
	if strand.Request().FormValue("direction") == "desc" {
		direction = "desc"
	}

	sql := strand.Request().FormValue("sql")
	sqlExecuted := sql
	params := []any{}
	if sql == "" {
		sql = "select * from " + trance.GetDialect().QuoteIdentifier(table) + " order by " + trance.GetDialect().QuoteIdentifier(order) + " " + direction
		sqlExecuted = sql + " limit ? offset ?"
		params = []any{describeTableLimit, offset}
	}
	res, err := trance.Database().Query(sqlExecuted, params...)
	if err != nil {
		return err
	}

	n, modals, rows, err := scanRows(fields, headers, res)
	if err != nil {
		return err
	}
	for i := range modals {
		modals[i]["table"] = table
	}

	var count int
	row := trance.Database().QueryRow("select count(*) from " + trance.GetDialect().QuoteIdentifier(table))
	if err := row.Scan(&count); err != nil {
		return err
	}
	next := ""
	if (offset + describeTableLimit) < count {
		next = fmt.Sprintf("/table/%s?order=%s&direction=%s&offset=%d", table, order, direction, offset+describeTableLimit)
	}
	prev := ""
	if offset > 0 {
		prev = fmt.Sprintf("/table/%s?order=%s&direction=%s&offset=%d", table, order, direction, max(0, offset-describeTableLimit))
	}

	links, err := describeSchemaSqlite()
	if err != nil {
		return err
	}
	for i := range links {
		if links[i]["label"] == table {
			links[i]["active"] = true
		}
	}

	return tableTemplate.Execute(strand.Response, map[string]any{
		"count":     count,
		"direction": direction,
		"end":       offset + n,
		"headers":   headers,
		"links":     links,
		"modals":    modals,
		"next":      next,
		"order":     order,
		"prev":      prev,
		"rows":      rows,
		"sql":       sql,
		"start":     offset + 1,
		"table":     table,
	})
}

func describeTableSqlite(table string) (map[string]reflect.StructField, []string, string, error) {
	fields := make(map[string]reflect.StructField, 0)
	headers := make([]string, 0)
	var pk string
	err := trance.Query[sqlitePragmaTableInfo]().
		SqlAll("select * from pragma_table_info(?) order by cid", table).
		Then(func(columns []*sqlitePragmaTableInfo) error {
			for _, column := range columns {
				if column.Pk {
					pk = column.Name
				}
				fields[column.Name] = reflect.StructField{
					Name: column.Name,
					Type: reflect.TypeOf(sql.NullString{}),
				}
				headers = append(headers, column.Name)
			}
			return nil
		}).
		Error

	return fields, headers, pk, err
}
