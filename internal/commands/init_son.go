package commands

import (
	"github.com/altiby/son/internal/common/imp"
	"github.com/altiby/son/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func CliImportUsersToDB(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "users",
		Usage: "import users from scv.zip file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "file",
				Usage:    "csv.zip file",
				Aliases:  []string{"f"},
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			i := imp.NewUsersImporter(ctx.String("file"), cfg)
			log.Info().Msgf("import users started")
			if err := i.Do(); err != nil {
				return cli.Exit(err.Error(), 1)
			}
			log.Info().Msgf("import users finished")
			return nil
		},
	}
}
