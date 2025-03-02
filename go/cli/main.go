package main

import (
	"encoding/json"
	"fmt"
	"os"

	"ale"

	"github.com/urfave/cli/v2"
)

// formatError wraps errors with context and ensures consistent error messages
func formatError(op string, err error) error {
	return fmt.Errorf("cli: %s: %w", op, err)
}

func main() {
	app := &cli.App{
		Name:  "golang-ale-cli",
		Usage: "CLI to read and write ALE files",
		Commands: []*cli.Command{
			{
				Name:  "read",
				Usage: "Read an ALE file",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "json",
						Usage:   "Output in JSON format",
						Aliases: []string{"j"},
					},
				},
				Action: func(c *cli.Context) error {
					// Validate input
					if c.NArg() < 1 {
						return formatError("read", fmt.Errorf("missing file path argument"))
					}
					inputFile := c.Args().Get(0)

					// Verify file exists and is readable
					if _, err := os.Stat(inputFile); err != nil {
						if os.IsNotExist(err) {
							return formatError("read", fmt.Errorf("file does not exist: %s", inputFile))
						}
						return formatError("read", fmt.Errorf("cannot access file: %s: %w", inputFile, err))
					}

					// Log input file
					fmt.Fprintf(c.App.Writer, "cli: Input file: %s\n", inputFile)

					// Create a new ALE handler
					handler := ale.New()

					// Read the ALE file
					aleObj, err := handler.ReadFile(inputFile)
					if err != nil {
						return formatError("read file", err)
					}

					// Output based on format
					if c.Bool("json") {
						// Marshal with indentation for readability
						jsonData, err := json.MarshalIndent(aleObj, "", "    ")
						if err != nil {
							return formatError("marshal JSON", err)
						}

						// Write JSON output
						c.App.Writer.Write(jsonData)
						c.App.Writer.Write([]byte("\n"))
					} else {
						// Write string representation
						fmt.Fprintf(c.App.Writer, "cli: Input ALE object: %s\n", aleObj.String())
					}
					return nil
				},
			},
			{
				Name:  "write",
				Usage: "Write metadata to an ALE file",
				Action: func(c *cli.Context) error {
					fmt.Fprintf(c.App.Writer, "cli: Write command executed\n")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
