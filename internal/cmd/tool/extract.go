package tool

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func (c *Cmd) extractCommand() *cli.Command {
	return &cli.Command{
		Name:  "extract",
		Usage: "Extract tools from archives to a temporary directory.",
		Commands: []*cli.Command{
			c.extractTarCommand(),
			c.extractTgzCommand(),
			c.extractZipCommand(),
		},
	}
}

func (c *Cmd) extractTarCommand() *cli.Command {
	return &cli.Command{
		Name:  "tar",
		Usage: "Extract a tar archive to a temporary directory.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Usage:    "Path to the tar archive.",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			path := cmd.String("path")

			slog.Debug("Extracting tar archive.", slog.String("path", path))

			p, err := c.ExtractTar(path, false)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Archive extracted successfully.")
			return nil
		},
	}
}

func (c *Cmd) extractTgzCommand() *cli.Command {
	return &cli.Command{
		Name:  "tgz",
		Usage: "Extract a tar.gz archive to a temporary directory.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Usage:    "Path to the tar.gz archive.",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			path := cmd.String("path")

			slog.Debug("Extracting tar.gz archive.", slog.String("path", path))

			p, err := c.ExtractTar(path, true)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Archive extracted successfully.")
			return nil
		},
	}
}

func (c *Cmd) extractZipCommand() *cli.Command {
	return &cli.Command{
		Name:  "zip",
		Usage: "Extract a zip archive to a temporary directory.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Usage:    "Path to the zip archive.",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			path := cmd.String("path")

			slog.Debug("Extracting zip archive.", slog.String("path", path))

			p, err := c.ExtractZip(path)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Archive extracted successfully.")
			return nil
		},
	}
}
