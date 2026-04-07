package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/action-stars/ghactl/internal/cmd/tool"
)

// This will be set by GoReleaser.
var version = "unknown"

func main() {
	app := &cli.Command{
		Name:    "ghactl",
		Usage:   "CLI to interact with GitHub Actions.",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose output.",
			},
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			level := slog.LevelWarn
			if cmd.Bool("verbose") {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
			return ctx, nil
		},
		Commands: []*cli.Command{
			tool.New(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("Fatal error.", slog.Any("error", err))
		os.Exit(1)
	}
}
