package tool

import (
	"context"
	"log/slog"

	"github.com/urfave/cli/v3"
)

func (c *Cmd) cacheCommand() *cli.Command {
	return &cli.Command{
		Name:  "cache",
		Usage: "Manage the tool cache.",
		Commands: []*cli.Command{
			c.cacheGetCommand(),
			c.cacheFindCommand(),
			c.cacheAddCommand(),
		},
	}
}

func (c *Cmd) cacheGetCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get the tool cache path.",
		Action: func(_ context.Context, cmd *cli.Command) error {
			slog.Debug("Getting the runner tool cache directory.")

			p, err := c.CacheGet()
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Runner tool cache directory retrieved.")
			return nil
		},
	}
}

func (c *Cmd) cacheFindCommand() *cli.Command {
	return &cli.Command{
		Name:  "find",
		Usage: "Find the path or versions of a tool.",
		Flags: []cli.Flag{
			toolNameFlag(),
			archFlag(),
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Return all matching cached versions.",
			},
			&cli.StringFlag{
				Name:  "version",
				Usage: "Version spec of the tool to find.",
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			tool := cmd.String("name")
			arch := defaultArch(cmd.String("arch"))
			all := cmd.Bool("all")
			versionSpec := cmd.String("version")
			if versionSpec == "" {
				versionSpec = "*"
			}

			slog.Debug("Finding tool.", slog.String("tool", tool), slog.String("versionSpec", versionSpec), slog.Bool("all", all))

			if all {
				vs, err := c.CacheFindAll(tool, arch)
				if err != nil {
					return exitErr(err)
				}

				if len(vs) == 0 {
					slog.Debug("Tool not found.", slog.String("tool", tool))
					return nil
				}

				count := 0

				for _, v := range vs {
					match, err := c.VersionCheck(v, versionSpec)
					if err != nil {
						return exitErr(err)
					}

					if !match {
						continue
					}

					if err := writeOutput(cmd, v); err != nil {
						return err
					}

					count++
				}

				if count == 0 {
					slog.Debug("Tool not found.", slog.String("tool", tool))
					return nil
				}

				slog.Debug("Tool versions found.", slog.String("tool", tool), slog.Int("versions", count))
				return nil
			}

			p, err := c.CacheFind(tool, arch, versionSpec)
			if err != nil {
				return exitErr(err)
			}

			if p == "" {
				slog.Debug("Tool not found.", slog.String("tool", tool))
				return nil
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Tool found.", slog.String("tool", tool))
			return nil
		},
	}
}

func (c *Cmd) cacheAddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Manage adding to the tool cache.",
		Commands: []*cli.Command{
			c.cacheAddDirCommand(),
			c.cacheAddFileCommand(),
		},
	}
}

func (c *Cmd) cacheAddDirCommand() *cli.Command {
	return &cli.Command{
		Name:  "dir",
		Usage: "Add a directory to the tool cache.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source",
				Usage:    "Source of the tool directory.",
				Required: true,
			},
			toolNameFlag(),
			versionFlag(),
			archFlag(),
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			source := cmd.String("source")
			tool := cmd.String("name")
			version := cmd.String("version")
			arch := defaultArch(cmd.String("arch"))

			slog.Debug("Adding tool to cache as directory.", slog.String("tool", tool), slog.String("version", version), slog.String("arch", arch))

			p, err := c.CacheDir(source, tool, version, arch)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Tool added to cache.", slog.String("tool", tool))
			return nil
		},
	}
}

func (c *Cmd) cacheAddFileCommand() *cli.Command {
	return &cli.Command{
		Name:  "file",
		Usage: "Add a file to the tool cache.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source",
				Usage:    "Source of the tool file.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "target-name",
				Usage: "Name to rename the source file to.",
			},
			toolNameFlag(),
			versionFlag(),
			archFlag(),
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			source := cmd.String("source")
			targetName := cmd.String("target-name")
			tool := cmd.String("name")
			version := cmd.String("version")
			arch := defaultArch(cmd.String("arch"))

			if targetName == "" {
				targetName = tool
			}

			slog.Debug("Adding tool to cache as file.", slog.String("tool", tool), slog.String("version", version), slog.String("arch", arch))

			p, err := c.CacheFile(source, targetName, tool, version, arch)
			if err != nil {
				return exitErr(err)
			}

			if err := writeOutput(cmd, p); err != nil {
				return err
			}

			slog.Debug("Tool added to cache.", slog.String("tool", tool))
			return nil
		},
	}
}
