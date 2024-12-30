package main

import (
	"github.com/alecthomas/kong"
	"github.com/evantbyrne/pocketsql"
)

var cli struct {
	Open pocketsql.OpenCommand `cmd:"" help:"Open database connection"`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
