package pocketsql

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/evantbyrne/trance"
)

//go:embed static
var static embed.FS

type OpenCommand struct {
	Connection string `arg:"" name:"connection" help:"Database connection string."`
	Port       string `flag:"" name:"port" default:"8080"`
}

func (cmd *OpenCommand) Run() error {
	return database(cmd.Connection, func() error {
		staticTree, err := fs.Sub(static, "static")
		if err != nil {
			return err
		}

		app := trance.App{ErrorHandler: errorHandler}
		app.Route("GET /", describeSchema)
		app.Route("GET /table/{name}", describeTable)
		app.Route("POST /query", customSql)
		http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticTree))))
		log.Println("Starting server on port", cmd.Port)
		return http.ListenAndServe(":"+cmd.Port, nil)
	})
}
