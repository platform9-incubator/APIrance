package main

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"./pkg/utils"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func loadAPISpec(c *cli.Context) error {
	swaggerDoc := c.Args().Get(0)
	// SwaggerDoc is a URL or path to a local file
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("All paths?? %s", specDoc.Analyzer.AllPaths())
	allPaths := specDoc.Analyzer.AllPaths()
	for s, pathItem := range allPaths {
		fmt.Printf("endpoint path: %s\n", s)
		fmt.Println(utils.Parse(&pathItem.PathItemProps))

	}
	//expanded, err := specDoc.Expanded()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(expanded)
	return nil
}

func main() {
	var baseCommandFlag string

	app := &cli.App{
		Name:    "APIrance",
		Usage:   "Makes your API make an appearance in Fission",
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
