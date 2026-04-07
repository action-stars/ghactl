package tool

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func (c *Cmd) versionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Manage tool versions.",
		Commands: []*cli.Command{
			c.versionCheckCommand(),
		},
	}
}

func (c *Cmd) versionCheckCommand() *cli.Command {
	return &cli.Command{
		Name:  "check",
		Usage: "Check if a version matches a constraint.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "version",
				Usage:    "Version to check.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "version-spec",
				Usage:    "Version spec to check against.",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			version := cmd.String("version")
			versionSpec := cmd.String("version-spec")

			slog.Debug("Checking version.", slog.String("version", version), slog.String("constraint", versionSpec))

			valid, err := c.VersionCheck(version, versionSpec)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, valid); err != nil {
				return err
			}

			slog.Debug("Version checked.")
			return nil
		},
	}
}
