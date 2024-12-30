package pocketsql

import (
	"github.com/evantbyrne/trance"
)

func describeSchema(strand *trance.Strand) error {
	links, err := describeSchemaSqlite()
	if err != nil {
		return err
	}
	if len(links) == 0 {
		return errorTemplate.Execute(strand.Response, map[string]any{
			"message": "Showing no tables",
			"links":   links,
		})
	}
	for _, link := range links {
		if link["label"] == "sqlite_schema" {
			return strand.Redirect("/table/sqlite_schema?order=name")
		}
	}
	// The `sqlite_schema` table should always exist, but we fall back to this just in case.
	return strand.Redirect(links[0]["url"].(string))
}

func describeSchemaSqlite() ([]map[string]any, error) {
	links := []map[string]any{}
	err := trance.Query[sqlitePragmaTableList]().
		SqlAll("select * from pragma_table_list order by name").
		Then(func(tables []*sqlitePragmaTableList) error {
			for _, table := range tables {
				links = append(links, map[string]any{
					"active": false,
					"label":  table.Name,
					"url":    "/table/" + table.Name,
				})
			}
			return nil
		}).
		Error

	return links, err
}
