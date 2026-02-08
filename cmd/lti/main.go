package main

import (
	"context"

	"github.com/alecthomas/kong"
	"github.com/go-playground/validator/v10"

	"github.com/julianstephens/liturgical-time-index/internal/commands"
	"github.com/julianstephens/liturgical-time-index/internal/util"
)

type CLI struct {
	Version  kong.VersionFlag     `short:"v" help:"Show version."`
	Build    commands.BuildCmd    `          help:"Build the index for a given year." cmd:"" name:"build"`
	Validate commands.ValidateCmd `          help:"Validate the plan file."           cmd:"" name:"validate"`
}

func main() {
	validate := validator.New()
	if err := validate.RegisterValidation("ISO8601", util.ValidateISO8601); err != nil {
		panic("Failed to register ISO8601 validation: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), util.GlobalTimeout)
	defer cancel()

	kongCtx := kong.Parse(
		&CLI{},
		kong.Name("lti"),
		kong.Description("CLI that compiles a Roman-season liturgical calendar into daily practice entries"),
		kong.Vars{"version": util.Version},
		kong.Bind(validate),
		kong.Bind(ctx),
	)

	err := kongCtx.Run()
	kongCtx.FatalIfErrorf(err)
}
