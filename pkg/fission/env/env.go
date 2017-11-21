package env

import (
	"github.com/go-openapi/spec"
	"os/exec"
	"fmt"
	"errors"
	"log"
	"os"
)

// Download remote file as another name
// curl -Lo test.js https://raw.githubusercontent.com/fission/fission/master/examples/nodejs/hello.js

// https://github.com/fission/fission/tree/master/examples/nodejs
const EXAMPLE_DIR string = "https://github.com/fission/fission/tree/master/examples/"

type FissionEnvironment struct {
	Name          string
	StubFileURL   string
	ImageName     string
	FileExtension string
}

var Fission = map[string]*FissionEnvironment{
	"nodejs": {
		Name:          "nodejs",
		ImageName:     "fission/node-env",
		StubFileURL:   "https://raw.githubusercontent.com/fission/fission/master/examples/nodejs/hello.js",
		FileExtension: ".js",
	},
	"python": {
		Name:          "python",
		ImageName:     "fission/python",
		StubFileURL:   "https://raw.githubusercontent.com/fission/fission/master/examples/python/hello.py",
		FileExtension: ".py",
	},
}

// fission env create --name nodejs --image fission/node-env
func (f FissionEnvironment) InitializeEnvironment() error {

	return nil
}

func (f FissionEnvironment) Run(httpMethod string, o *spec.Operation) error {
	f.createFunction(httpMethod, o)
	f.createRoute(httpMethod, o.ID, o)
	return nil
}

// fission function create --name hello --env nodejs --code hello.js
func (f FissionEnvironment) createFunction(httpMethod string, o *spec.Operation) error {
	// Need to get example function for this environment and rename it before creating the fission function
	_, created, downloadErr := f.downloadStub(httpMethod, o.ID)
	if downloadErr != nil {
		log.Fatalln(downloadErr)
		return downloadErr
	} else if created == false {
		// Function file with same name in same path already exists.
		return nil
	}

	return nil
}

// DownloadStub downloads the FissionEnvironment's stub file but renamed to funcName
// httpDir is one of GET|POST|PUT|DELETE|HEAD which is a dir in the root of the project
func (f FissionEnvironment) downloadStub(httpMethod, funcName string) (pathToFile string, created bool, err error) {
	functionFile := fmt.Sprintf("%s/%s%s", httpMethod, funcName, f.FileExtension)
	if _, err := os.Stat(functionFile); err == nil {
		return "", false, nil
	}
	curl := exec.Command("curl", "-Lo", functionFile, f.StubFileURL)
	_, err = curl.Output()
	if err != nil {
		fmt.Println(err)
		return "", false, errors.New(fmt.Sprintf("error downloading stub file %s", f.StubFileURL))
	}
	return pathToFile, true, err
}

// fission route create --method GET --url /hello --function hello
func (f FissionEnvironment) createRoute(httpMethod, funcName string, o *spec.Operation) error {
	f.downloadStub(httpMethod, o.ID)
	return nil
}

