package main

import (
	"github.com/alecthomas/kong"
	"github.com/go-playground/validator/v10"

	"github.com/julianstephens/liturgical-time-index/internal/commands"
	"github.com/julianstephens/liturgical-time-index/internal/util"
)

type CLI struct {
	Version kong.VersionFlag  `short:"v" help:"Show version."`
	Build   commands.BuildCmd `          help:"Build the index for a given year." cmd:"" name:"build"`
}

func main() {
	validate := validator.New()
	if err := validate.RegisterValidation("ISO8601", util.ValidateISO8601); err != nil {
		panic("Failed to register ISO8601 validation: " + err.Error())
	}

	kongCtx := kong.Parse(
		&CLI{},
		kong.Name("lti"),
		kong.Description("CLI that compiles a Roman-season liturgical calendar into daily practice entries"),
		kong.Vars{"version": "0.1.0"},
		kong.Bind(validate),
	)

	err := kongCtx.Run()
	kongCtx.FatalIfErrorf(err)
}
