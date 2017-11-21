package env

import (
	"github.com/go-openapi/spec"
	"os/exec"
	"fmt"
	"errors"
	"log"
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
	_, downloadErr := f.downloadStub(httpMethod, o.ID)
	if downloadErr != nil {
		log.Fatalln(downloadErr)
		return downloadErr
	}

	return nil
}

// DownloadStub downloads the FissionEnvironment's stub file but renamed to funcName
// httpDir is one of GET|POST|PUT|DELETE|HEAD which is a dir in the root of the project
func (f FissionEnvironment) downloadStub(httpMethod, funcName string) (pathToFile string, err error) {
	functionFile := fmt.Sprintf("%s/%s%s", httpMethod, funcName, f.FileExtension)
	curl := exec.Command("curl", "-Lo", functionFile, f.StubFileURL)
	_, err = curl.Output()
	if err != nil {
		fmt.Println(err)
		return "", errors.New(fmt.Sprintf("error downloading stub file %s", f.StubFileURL))
	}
	return pathToFile, err
}

// fission route create --method GET --url /hello --function hello
func (f FissionEnvironment) createRoute(httpMethod, funcName string, o *spec.Operation) error {
	f.downloadStub(httpMethod, o.ID)
	return nil
}

