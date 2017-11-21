package main

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"./pkg/fission/env"
	"./pkg/fission"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"strings"
)

// CWD doesn't have trailing slash
var CWD string
var FISSION_BIN string


func init() {
	// Ensure we can determine our current working directory
	CWD, err := os.Getwd()
	if err != nil {
		log.Println("Literally how man??")
		log.Fatalln(err)
	}
	// TODO: delete this print statement
	fmt.Printf("Current working directory is %s\n", CWD)
	// Check if fission CLI is available to use

	which := exec.Command("which", "fission")
	FISSION_BIN, err := which.Output()
	if err != nil {
		log.Println(`
		Fission tool is not installed. Assuming you're on MacOS, run below command
		curl -Lo fission https://github.com/fission/fission/releases/download/v0.4.0/fission-cli-osx
		&& chmod +x fission
		&& sudo mv fission /usr/local/bin/
		`)
		log.Fatalln(err)
	}
	// TODO: delete this print statement
	fmt.Printf("Found fission at %s\n", FISSION_BIN)
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

func generate(c *cli.Context) error {
	var fizzshin *env.FissionEnvironment
	var swaggerDoc string
	// Ensure fission environment flag value and path/URL for swagger spec is provided
	if fissEnv := c.String("env"); fissEnv == "" || c.NArg() <= 0 {
		cli.ShowAppHelp(c)
		log.Fatalln("not enough arguments")
	} else {
		fizzshin = env.Fission[fissEnv]
		fizzshin.InitializeEnvironment()
		swaggerDoc = c.Args().First()
	}

	// swaggerDoc is a URL or path to a local file
	specDoc, err := loads.Spec(swaggerDoc)
	if err != nil {
		log.Fatalln(err)
	}
	scaffoldAPI(fizzshin, specDoc)
	return nil
}

func scaffoldAPI(f *env.FissionEnvironment ,d *loads.Document) {
	var wg sync.WaitGroup
	operations := d.Analyzer.Operations()
	for httpMethod, v := range operations {
		scaffoldHTTPMethodDir(httpMethod)

		for _, operation := range v {
			// FIXME: below hack will ignore endpoints like '/pets/{id}' due to
			// the operation ID being a sentence
			if len(strings.Split(operation.ID, " ")) > 1 {
				continue
			}
			wg.Add(1)
			go fission.Start(&fission.Worker{
				WaitGroup: &wg,
				HttpMethod: httpMethod,
				Operation: operation,
				Environment: f,
			})
		}
	}
	wg.Wait()
}

// scaffoldHTTPMethodDir looks in the current dir for a dir matching httpMethod
// and creates it if not there
func scaffoldHTTPMethodDir(httpMethod string) error {
	methodDir := filepath.Join(CWD, httpMethod)
	err := os.MkdirAll(methodDir, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func main() {
	app := &cli.App{
		Name:    "APIrance",
		Usage:   "Makes your API make an appearance in Fission",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Usage: "Environment and language in which routes and functions will be created",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate some yaml files from this spec",
				Action:  generate,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "commit"},
					&cli.BoolFlag{Name: "push"},
				},
			},
		},
	}

	app.Run(os.Args)
}
