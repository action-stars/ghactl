package tool

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

func (c *Cmd) installCommand() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "Install and cache a tool release.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "owner",
				Usage:    "GitHub repository owner.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repo",
				Usage:    "GitHub repository name.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "Tool cache name. Defaults to --repo.",
			},
			&cli.StringFlag{
				Name:  "version",
				Usage: "Version input. Supports latest or exact version/tag.",
				Value: "latest",
			},
			archFlag(),
			osFlag(),
			&cli.BoolFlag{
				Name:  "pre-release",
				Usage: "Include pre-releases when resolving latest.",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "GitHub token. Defaults to GITHUB_TOKEN when unset.",
			},
			&cli.BoolFlag{
				Name:  "add-to-path",
				Usage: "Add the tool directory to PATH.",
				Value: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			owner := cmd.String("owner")
			repo := cmd.String("repo")
			name := cmd.String("name")
			version := strings.TrimSpace(cmd.String("version"))
			arch := defaultArch(cmd.String("arch"))
			osName := cmd.String("os")
			includePreRelease := cmd.Bool("pre-release")
			token := cmd.String("token")
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			addToPath := cmd.Bool("add-to-path")

			slog.Debug("Installing tool.",
				slog.String("owner", owner),
				slog.String("repo", repo),
				slog.String("name", name),
				slog.String("version", version),
				slog.String("arch", arch),
				slog.String("os", osName),
				slog.Bool("preRelease", includePreRelease),
				slog.Bool("addToPath", addToPath),
			)

			p, err := c.Install(ctx, InstallOptions{
				Name:              name,
				Owner:             owner,
				Repo:              repo,
				Version:           version,
				Arch:              arch,
				OS:                osName,
				IncludePreRelease: includePreRelease,
				Token:             token,
				AddToPath:         addToPath,
			})
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Tool installed.", slog.String("path", p))
			return nil
		},
	}
}
