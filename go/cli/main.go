package main

import (
	"os"

	"lib-post-interchange/libale"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "libale-go-cli",
		Usage: "A CLI to read and write ALE files",
		Commands: []*cli.Command{
			{
				Name:  "read",
				Usage: "Read an ALE file",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return cli.Exit("please provide a file path as the first argument", 1)
					}
					inputFile := c.Args().Get(0)
					c.App.Writer.Write([]byte("cli: Input file: " + inputFile + "\n"))

					// Create a new ALE handler
					handler := libale.New()

					// Read the ALE file
					ale, err := handler.ReadFile(inputFile)
					if err != nil {
						return cli.Exit("libale: error reading file: "+err.Error(), 1)
					}
					c.App.Writer.Write([]byte("cli: Output ALE: " + ale.String() + "\n"))
					return nil
				},
			},
			{
				Name:  "write",
				Usage: "Write metadata to an ALE file",
				Action: func(c *cli.Context) error {
					// Function stub for writing metadata
					c.App.Writer.Write([]byte("cli: Write command executed\n"))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		cli.OsExiter(1)
	}
}
