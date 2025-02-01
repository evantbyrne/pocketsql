package pocketsql

import (
	"database/sql"

	"github.com/evantbyrne/trance"
)

type sqlitePragmaTableInfo struct {
	Cid       string         `@:"cid"`
	Name      string         `@:"name"`
	Type      string         `@:"type"`
	NotNull   int            `@:"notnull"`
	TableName string         `@:"tbl_name"`
	Sql       string         `@:"sql"`
	DfltValue sql.NullString `@:"dflt_value"`
	Pk        bool           `@:"pk"`
}

type sqlitePragmaTableList struct {
	Schema string `@:"schema"`
	Name   string `@:"name"`
	Type   string `@:"type"`
	NCol   int    `@:"ncol"`
	Wr     int    `@:"wr"`
	Strict int    `@:"strict"`
}

func (sqlitePragmaTableList) WeaveConfig() trance.WeaveConfig {
	return trance.WeaveConfig{Table: "pragma_table_list"}
}
