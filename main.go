package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"

	"github.com/action-stars/ghactl/internal/toolkit/core"
	"github.com/action-stars/ghactl/internal/toolkit/toolcache"
	"github.com/carlmjohnson/versioninfo"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "ghactl",
		Usage:   "CLI to interact with GitHub Actions",
		Version: versioninfo.Version,
		Commands: []*cli.Command{
			{
				Name:  "env",
				Usage: "Manage workflow environment variables",
				Subcommands: []*cli.Command{
					{
						Name:  "clear",
						Usage: "Clear an environment variable",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "key",
								Usage:    "Environment variable key",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							key := cCtx.String("key")
							err := core.ExportVariable(key, "")
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err = fmt.Fprintf(cCtx.App.ErrWriter, "Environment variable %s cleared\n", key)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
					{
						Name:  "set",
						Usage: "Set an environment variable",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "key",
								Usage:    "Environment variable key",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Environment variable value",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							key := cCtx.String("key")
							err := core.ExportVariable(key, cCtx.String("value"))
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Environment variable %s set\n", key)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "matcher",
				Usage: "Manage problem matchers",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add a problem matcher",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "config",
								Usage:    "Config file path",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							p := cCtx.String("config")

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Adding problem matcher from %s\n", p)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}

							err := core.AddMatcher(cCtx.App.Writer, p)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Problem matcher added from %s\n", p)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "Remove a problem matcher",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "owner",
								Usage:    "Owner of the problem matcher",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							owner := cCtx.String("owner")

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Removing problem matcher with owner %s\n", owner)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}

							err := core.RemoveMatcher(cCtx.App.Writer, owner)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Problem matcher for owner %s removed\n", owner)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "message",
				Usage: "Manage workflow log messages",
				Subcommands: []*cli.Command{
					{
						Name:  "group",
						Usage: "Manage log grouping",
						Subcommands: []*cli.Command{
							{
								Name:  "start",
								Usage: "Start a log group",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "name",
										Usage:    "Log group name",
										Required: true,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									err := core.StartGroup(cCtx.App.Writer, cCtx.String("name"))
									if err != nil {
										return cli.Exit(err, 1)
									}
									return nil
								},
							},
							{
								Name:  "end",
								Usage: "End a log group",
								Flags: []cli.Flag{
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									err := core.EndGroup(cCtx.App.Writer)
									if err != nil {
										return cli.Exit(err, 1)
									}
									return nil
								},
							},
						},
					},
					{
						Name:  "debug",
						Usage: "Write a debug log message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "message",
								Usage:    "Log message to write",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							err := core.Debug(cCtx.App.Writer, cCtx.String("message"))
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					{
						Name:  "info",
						Usage: "Write an info log message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "message",
								Usage:    "Log message to write",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "Custom title for the log annotation",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "file",
								Usage:    "File where the log annotation should be added",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "column",
								Usage:    "Column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-column",
								Usage:    "End column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "line",
								Usage:    "Line in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-line",
								Usage:    "End line in the file to add the log annotation",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							column := -1
							if cCtx.IsSet("column") {
								column = cCtx.Int("column")
							}
							endColumn := -1
							if cCtx.IsSet("end-column") {
								endColumn = cCtx.Int("end-column")
							}
							line := -1
							if cCtx.IsSet("line") {
								line = cCtx.Int("line")
							}
							endLine := -1
							if cCtx.IsSet("end-line") {
								endLine = cCtx.Int("end-line")
							}

							err := core.Info(cCtx.App.Writer, cCtx.String("message"), core.AnnotationProperties{
								Title:     cCtx.String("title"),
								File:      cCtx.String("file"),
								Column:    column,
								EndColumn: endColumn,
								Line:      line,
								EndLine:   endLine,
							})
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					{
						Name:  "notice",
						Usage: "Write a notice log message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "message",
								Usage:    "Log message to write",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "Custom title for the log annotation",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "file",
								Usage:    "File where the log annotation should be added",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "column",
								Usage:    "Column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-column",
								Usage:    "End column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "line",
								Usage:    "Line in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-line",
								Usage:    "End line in the file to add the log annotation",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							column := -1
							if cCtx.IsSet("column") {
								column = cCtx.Int("column")
							}
							endColumn := -1
							if cCtx.IsSet("end-column") {
								endColumn = cCtx.Int("end-column")
							}
							line := -1
							if cCtx.IsSet("line") {
								line = cCtx.Int("line")
							}
							endLine := -1
							if cCtx.IsSet("end-line") {
								endLine = cCtx.Int("end-line")
							}

							err := core.Notice(cCtx.App.Writer, cCtx.String("message"), core.AnnotationProperties{
								Title:     cCtx.String("title"),
								File:      cCtx.String("file"),
								Column:    column,
								EndColumn: endColumn,
								Line:      line,
								EndLine:   endLine,
							})
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					{
						Name:  "warn",
						Usage: "Write a warning log message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "message",
								Usage:    "Log message to write",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "Custom title for the log annotation",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "file",
								Usage:    "File where the log annotation should be added",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "column",
								Usage:    "Column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-column",
								Usage:    "End column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "line",
								Usage:    "Line in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-line",
								Usage:    "End line in the file to add the log annotation",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							column := -1
							if cCtx.IsSet("column") {
								column = cCtx.Int("column")
							}
							endColumn := -1
							if cCtx.IsSet("end-column") {
								endColumn = cCtx.Int("end-column")
							}
							line := -1
							if cCtx.IsSet("line") {
								line = cCtx.Int("line")
							}
							endLine := -1
							if cCtx.IsSet("end-line") {
								endLine = cCtx.Int("end-line")
							}

							err := core.Warning(cCtx.App.Writer, cCtx.String("message"), core.AnnotationProperties{
								Title:     cCtx.String("title"),
								File:      cCtx.String("file"),
								Column:    column,
								EndColumn: endColumn,
								Line:      line,
								EndLine:   endLine,
							})
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
					{
						Name:  "error",
						Usage: "Write an error log message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "message",
								Usage:    "Log message to write",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "title",
								Usage:    "Custom title for the log annotation",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "file",
								Usage:    "File where the log annotation should be added",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "column",
								Usage:    "Column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-column",
								Usage:    "End column in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "line",
								Usage:    "Line in the file to add the log annotation",
								Required: false,
							},
							&cli.IntFlag{
								Name:     "end-line",
								Usage:    "End line in the file to add the log annotation",
								Required: false,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							column := -1
							if cCtx.IsSet("column") {
								column = cCtx.Int("column")
							}
							endColumn := -1
							if cCtx.IsSet("end-column") {
								endColumn = cCtx.Int("end-column")
							}
							line := -1
							if cCtx.IsSet("line") {
								line = cCtx.Int("line")
							}
							endLine := -1
							if cCtx.IsSet("end-line") {
								endLine = cCtx.Int("end-line")
							}

							err := core.Error(cCtx.App.Writer, cCtx.String("message"), core.AnnotationProperties{
								Title:     cCtx.String("title"),
								File:      cCtx.String("file"),
								Column:    column,
								EndColumn: endColumn,
								Line:      line,
								EndLine:   endLine,
							})
							if err != nil {
								return cli.Exit(err, 1)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "output",
				Usage: "Manage workflow step outputs",
				Subcommands: []*cli.Command{
					{
						Name:  "set",
						Usage: "Set an workflow step output",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "key",
								Usage:    "Output key",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Output value",
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "verbose",
								Usage: "Enable verbose output",
							},
						},
						Action: func(cCtx *cli.Context) error {
							key := cCtx.String("key")
							value := cCtx.String("value")

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Setting output value %s to %s\n", key, value)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}

							err := core.SetOutput(key, value)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Output value %s set\n", key)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "path",
				Usage: "Manage workflow PATH configuration",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add an entry to the workflow PATH",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Path value",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							value := cCtx.String("value")
							err := core.AddPath(value)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Value %s added to PATH\n", value)
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "secret",
				Usage: "Manage workflow secrets",
				Subcommands: []*cli.Command{
					{
						Name:  "mask",
						Usage: "Mask a workflow secret value",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Value to mask",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							err := core.SetSecret(cCtx.App.Writer, cCtx.String("value"))
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Secret value masked")
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
					{
						Name:  "encrypt",
						Usage: "Encrypt a workflow secret value",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret",
								Usage:    "Encryption secret",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Value to encrypt",
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "verbose",
								Usage: "Enable verbose output",
							},
						},
						Action: func(cCtx *cli.Context) error {
							cipherText, err := core.EncryptSecret(cCtx.String("secret"), cCtx.String("value"))
							if err != nil {
								return cli.Exit(err, 1)
							}

							_, err = fmt.Fprintln(cCtx.App.Writer, cipherText)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Secret value encrypted")
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
					{
						Name:  "decrypt",
						Usage: "Decrypt a workflow secret value",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "secret",
								Usage:    "Encryption secret",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "secret-value",
								Usage:    "Secret value to decrypt",
								Required: true,
							},
							&cli.BoolFlag{
								Name:  "verbose",
								Usage: "Enable verbose output",
							},
						},
						Action: func(cCtx *cli.Context) error {
							plainText, err := core.DecryptSecret(cCtx.String("secret"), cCtx.String("secret-value"))
							if err != nil {
								return cli.Exit(err, 1)
							}

							_, err = fmt.Fprintln(cCtx.App.Writer, plainText)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Secret value decrypted")
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "summary",
				Usage: "Manage workflow summary messages",
				Subcommands: []*cli.Command{
					{
						Name:  "write",
						Usage: "Write a workflow step summary message",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "value",
								Usage:    "Path value",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							value := cCtx.String("value")
							err := core.WriteSummary(value)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if cCtx.Bool("verbose") {
								_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Step summary written")
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "temp",
				Usage: "Get the runner temporary directory",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "verbose",
						Usage: "Enable verbose output",
					},
				},
				Action: func(cCtx *cli.Context) error {
					verbose := cCtx.Bool("verbose")

					if verbose {
						_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Getting the runner temporary directory")
						if err != nil {
							return cli.Exit(err, 1)
						}
					}

					d, err := core.GetTempDirectory()
					if err != nil {
						return cli.Exit(err, 1)
					}

					_, err = fmt.Fprintln(cCtx.App.Writer, d)
					if err != nil {
						return cli.Exit(err, 1)
					}

					if verbose {
						_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Runner temporary directory retrieved")
						if err != nil {
							return cli.Exit(err, 1)
						}
					}
					return nil
				},
			},
			{
				Name:  "tool",
				Usage: "Manage GitHub runner tools",
				Subcommands: []*cli.Command{
					{
						Name:  "cache",
						Usage: "Manage the tool cache",
						Subcommands: []*cli.Command{
							{
								Name:  "get",
								Usage: "Gets the tool cache path",
								Flags: []cli.Flag{
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Getting the runner tool cache directory")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									p, err := toolcache.GetToolCacheDirectory()
									if err != nil {
										return cli.Exit(err, 1)
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, p)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Runner tool cache directory retrieved")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							{
								Name:  "find-all",
								Usage: "Find all the paths to a tool",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "name",
										Usage:    "Name of the tool to find",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "arch",
										Usage:    "Architecture of the tool to find",
										Required: false,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									tool := cCtx.String("name")
									arch := cCtx.String("arch")
									if arch == "" {
										arch = runtime.GOARCH
									}

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Finding all paths for tool %s\n", tool)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									ps, err := toolcache.FindAllToolVersions(tool, arch)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if len(ps) == 0 {
										if verbose {
											_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s not found\n", tool)
											if err != nil {
												return cli.Exit(err, 1)
											}
										}
										return nil
									}

									for _, p := range ps {
										_, err := fmt.Fprintln(cCtx.App.Writer, p)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s found at %d paths\n", tool, len(ps))
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							{
								Name:  "find",
								Usage: "Find the path to a tool",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "name",
										Usage:    "Name of the tool to find",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "arch",
										Usage:    "Architecture of the tool to find",
										Required: false,
									},
									&cli.StringFlag{
										Name:     "version",
										Usage:    "Version spec of the tool to find",
										Required: false,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									tool := cCtx.String("name")
									arch := cCtx.String("arch")
									if arch == "" {
										arch = runtime.GOARCH
									}
									versionSpec := cCtx.String("version")
									if versionSpec == "" {
										versionSpec = "*"
									}

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Finding tool %s with version spec %s\n", tool, versionSpec)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									p, err := toolcache.FindTool(tool, arch, versionSpec)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if len(p) == 0 {
										if verbose {
											_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s not found\n", tool)
											if err != nil {
												return cli.Exit(err, 1)
											}
										}
										return nil
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, p)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s found\n", tool)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							{
								Name:  "add",
								Usage: "Manage adding to the tool cache",
								Subcommands: []*cli.Command{
									{
										Name:  "dir",
										Usage: "Add a directory to the tool cache",
										Flags: []cli.Flag{
											&cli.StringFlag{
												Name:     "source",
												Usage:    "Source of the tool directory",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "name",
												Usage:    "Name of the tool to add",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "version",
												Usage:    "Version of the tool to add",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "arch",
												Usage:    "Architecture of the tool to add",
												Required: false,
											},
											&cli.BoolFlag{
												Name:     "verbose",
												Usage:    "Enable verbose output",
												Required: false,
											},
										},
										Action: func(cCtx *cli.Context) error {
											verbose := cCtx.Bool("verbose")
											source := cCtx.String("source")
											tool := cCtx.String("name")
											version := cCtx.String("version")
											arch := cCtx.String("arch")
											if arch == "" {
												arch = runtime.GOARCH
											}

											if verbose {
												_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Adding tool %s version %s arch %s to the cache as a directory\n", tool, version, arch)
												if err != nil {
													return cli.Exit(err, 1)
												}
											}

											p, err := toolcache.CacheDir(source, tool, version, arch)
											if err != nil {
												return cli.Exit(err, 1)
											}

											_, err = fmt.Fprintln(cCtx.App.Writer, p)
											if err != nil {
												return cli.Exit(err, 1)
											}

											if verbose {
												_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s added to cache\n", tool)
												if err != nil {
													return cli.Exit(err, 1)
												}
											}
											return nil
										},
									},
									{
										Name:  "file",
										Usage: "Add a file to the tool cache",
										Flags: []cli.Flag{
											&cli.StringFlag{
												Name:     "source",
												Usage:    "Source of the tool directory",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "target-name",
												Usage:    "Name to rename the source file to",
												Required: false,
											},
											&cli.StringFlag{
												Name:     "name",
												Usage:    "Name of the tool to add",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "version",
												Usage:    "Version of the tool to add",
												Required: true,
											},
											&cli.StringFlag{
												Name:     "arch",
												Usage:    "Architecture of the tool to add",
												Required: false,
											},
											&cli.BoolFlag{
												Name:     "verbose",
												Usage:    "Enable verbose output",
												Required: false,
											},
										},
										Action: func(cCtx *cli.Context) error {
											verbose := cCtx.Bool("verbose")
											source := cCtx.String("source")
											targetName := cCtx.String("target-name")
											tool := cCtx.String("name")
											if targetName == "" {
												targetName = tool
											}
											version := cCtx.String("version")
											arch := cCtx.String("arch")
											if arch == "" {
												arch = runtime.GOARCH
											}

											if verbose {
												_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Adding tool %s version %s arch %s to the cache as a file\n", tool, version, arch)
												if err != nil {
													return cli.Exit(err, 1)
												}
											}

											p, err := toolcache.CacheFile(source, targetName, tool, version, arch)
											if err != nil {
												return cli.Exit(err, 1)
											}

											_, err = fmt.Fprintln(cCtx.App.Writer, p)
											if err != nil {
												return cli.Exit(err, 1)
											}

											if verbose {
												_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Tool %s added to cache\n", tool)
												if err != nil {
													return cli.Exit(err, 1)
												}
											}
											return nil
										},
									},
								},
							},
						},
					},
					{
						Name:  "download",
						Usage: "Download a tool to a temporary directory",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "url",
								Usage:    "URL to download the tool from",
								Required: true,
							},
							&cli.BoolFlag{
								Name:     "verbose",
								Usage:    "Enable verbose output",
								Required: false,
							},
						},
						Action: func(cCtx *cli.Context) error {
							verbose := cCtx.Bool("verbose")
							url, err := url.Parse(cCtx.String("url"))
							if err != nil {
								return cli.Exit(err, 1)
							}

							if verbose {
								_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Downloading tool from %s\n", url.String())
								if err != nil {
									return cli.Exit(err, 1)
								}
							}

							p, err := toolcache.DownloadTool(url.String())
							if err != nil {
								return cli.Exit(err, 1)
							}

							_, err = fmt.Fprintln(cCtx.App.Writer, p)
							if err != nil {
								return cli.Exit(err, 1)
							}

							if verbose {
								_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Tool downloaded successfully")
								if err != nil {
									return cli.Exit(err, 1)
								}
							}
							return nil
						},
					},
					{
						Name:  "extract",
						Usage: "Extract tools from archives to a temporary directory",
						Subcommands: []*cli.Command{
							{
								Name:  "tar",
								Usage: "Extract tar archive to a temporary directory",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "path",
										Usage:    "Path to the tar archive",
										Required: true,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									path := cCtx.String("path")

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Extracting tar archive at %s\n", path)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									p, err := toolcache.ExtractTar(path, false)
									if err != nil {
										return cli.Exit(err, 1)
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, p)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Archive extracted successfully")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							{
								Name:  "tgz",
								Usage: "Extract tar.gz archive to a temporary directory",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "path",
										Usage:    "Path to the tar.gz archive",
										Required: true,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									path := cCtx.String("path")

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Extracting tar.gz archive at %s\n", path)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									p, err := toolcache.ExtractTar(path, true)
									if err != nil {
										return cli.Exit(err, 1)
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, p)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Archive extracted successfully")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							{
								Name:  "zip",
								Usage: "Extract zip archive to a temporary directory",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "path",
										Usage:    "Path to the zip archive",
										Required: true,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									path := cCtx.String("path")

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Extracting zip archive at %s\n", path)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									p, err := toolcache.ExtractZip(path)
									if err != nil {
										return cli.Exit(err, 1)
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, p)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Archive extracted successfully")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
						},
					},
					{
						Name:  "version",
						Usage: "Manage tool versions",
						Subcommands: []*cli.Command{
							{
								Name:  "check",
								Usage: "Check if a version matches a constraint",
								Flags: []cli.Flag{
									&cli.StringFlag{
										Name:     "version",
										Usage:    "Version to check",
										Required: true,
									},
									&cli.StringFlag{
										Name:     "version-spec",
										Usage:    "Version spec to check against",
										Required: true,
									},
									&cli.BoolFlag{
										Name:     "verbose",
										Usage:    "Enable verbose output",
										Required: false,
									},
								},
								Action: func(cCtx *cli.Context) error {
									verbose := cCtx.Bool("verbose")
									version := cCtx.String("version")
									versionSpec := cCtx.String("version-spec")

									if verbose {
										_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Checking version %s against constraint %s\n", version, versionSpec)
										if err != nil {
											return cli.Exit(err, 1)
										}
									}

									valid, err := toolcache.CheckVersion(version, versionSpec)
									if err != nil {
										return cli.Exit(err, 1)
									}

									_, err = fmt.Fprintln(cCtx.App.Writer, valid)
									if err != nil {
										return cli.Exit(err, 1)
									}

									if verbose {
										_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Version checked")
										if err != nil {
											return cli.Exit(err, 1)
										}
									}
									return nil
								},
							},
							// {
							// 	Name:  "lookup",
							// 	Usage: "Lookup a version from GitHub",
							// 	Flags: []cli.Flag{
							// 		&cli.StringFlag{
							// 			Name:     "repository",
							// 			Usage:    "Repository to lookup the version from",
							// 			Required: true,
							// 		},
							// 		&cli.StringFlag{
							// 			Name:     "version",
							// 			Usage:    "Version spec to check against",
							// 			Required: true,
							// 		},
							// 		&cli.BoolFlag{
							// 			Name:     "pre-releases",
							// 			Usage:    "Include pre-releases",
							// 			Required: false,
							// 		},
							// 		&cli.BoolFlag{
							// 			Name:     "verbose",
							// 			Usage:    "Enable verbose output",
							// 			Required: false,
							// 		},
							// 	},
							// 	Action: func(cCtx *cli.Context) error {
							// 		verbose := cCtx.Bool("verbose")
							// 		repository := cCtx.String("repository")
							// 		versionSpec := cCtx.String("version")
							// 		preReleases := cCtx.Bool("pre-releases")

							// 		if verbose {
							// 			_, err := fmt.Fprintf(cCtx.App.ErrWriter, "Looking up version from GitHub repository %s with constraint %s\n", repository, versionSpec)
							// 			if err != nil {
							// 				return cli.Exit(err, 1)
							// 			}
							// 		}

							// 		v, err := toolkit.LookupGitHubVersion(repository, versionSpec, preReleases)
							// 		if err != nil {
							// 			return cli.Exit(err, 1)
							// 		}

							// 		if len(v) == 0 {
							// 			if verbose {
							// 				_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Matching version not found")
							// 				if err != nil {
							// 					return cli.Exit(err, 1)
							// 				}
							// 			}
							// 			return nil
							// 		}

							// 		_, err = fmt.Fprintln(cCtx.App.Writer, v)
							// 		if err != nil {
							// 			return cli.Exit(err, 1)
							// 		}

							// 		if verbose {
							// 			_, err := fmt.Fprintln(cCtx.App.ErrWriter, "Version looked up")
							// 			if err != nil {
							// 				return cli.Exit(err, 1)
							// 			}
							// 		}
							// 		return nil
							// 	},
							// },
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
