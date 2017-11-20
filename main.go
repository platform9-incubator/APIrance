package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"os"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"

	"log"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func loadAPISpec(c *cli.Context) error {
	swaggerDoc := c.Args().Get(0)
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(specDoc.Spec().SwaggerProps.Definitions)
	return nil
}

func main() {
	var baseCommandFlag string

	app := &cli.App{
		Name:  "APIrance",
		Usage: "Makes your API make an appearance in Fission",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "PH",
				Value:       "Placeholder",
				Usage:       "Used for placeholding a spot for adding flags to the base command",
				Destination: &baseCommandFlag,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate some yaml files from this spec",
				Action:  loadAPISpec,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "commit"},
					&cli.BoolFlag{Name: "push"},
				},
			},
		},
	}

	app.Run(os.Args)
}
