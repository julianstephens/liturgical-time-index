package main

import (
	"context"

	"github.com/alecthomas/kong"

	"github.com/julianstephens/liturgical-time-index/internal"
	"github.com/julianstephens/liturgical-time-index/internal/command"
)

type CLI struct {
	Version  kong.VersionFlag    `short:"v" help:"Show version."`
	Build    command.BuildCmd    `          help:"Build the index for a given year."  cmd:"" name:"build"`
	Today    command.TodayCmd    `          help:"Get the entry for a specific date." cmd:"" name:"today"`
	Validate command.ValidateCmd `          help:"Validate the plan file."            cmd:"" name:"validate"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), internal.GlobalTimeout)
	defer cancel()

	kongCtx := kong.Parse(
		&CLI{},
		kong.Name("lti"),
		kong.Description("CLI that compiles a Roman-season liturgical calendar into daily practice entries"),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{"version": internal.Version},
		kong.Bind(ctx),
	)

	err := kongCtx.Run()
	kongCtx.FatalIfErrorf(err)
}
