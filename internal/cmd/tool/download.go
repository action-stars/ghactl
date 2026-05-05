package tool

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func (c *Cmd) downloadCommand() *cli.Command {
	return &cli.Command{
		Name:  "download",
		Usage: "Download a tool to a temporary directory.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "URL to download the tool from.",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			rawURL := cmd.String("url")

			slog.Debug("Downloading tool.", slog.String("url", rawURL))

			p, err := c.Download(ctx, rawURL)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Tool downloaded successfully.")
			return nil
		},
	}
}
