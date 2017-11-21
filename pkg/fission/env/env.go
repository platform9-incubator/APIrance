package env

import (
	"github.com/go-openapi/spec"
	"os/exec"
	"fmt"
	"errors"
	"log"
	"os"
	"strings"
)

// Download remote file as another name
// curl -Lo test.js https://raw.githubusercontent.com/fission/fission/master/examples/nodejs/hello.js

// https://github.com/fission/fission/tree/master/examples/nodejs
const EXAMPLE_DIR string = "https://github.com/fission/fission/tree/master/examples/"
var FISSION_BIN string = getFissionBinary()

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

func runFission(command string) string {
	cmd := fmt.Sprintf("%s %s", getFissionBinary(), command)
	fmt.Println("cmd: "+cmd)
	listFuncs := exec.Command("/usr/bin/bash", "-c", cmd)
	var out []byte
	out, err := listFuncs.CombinedOutput()
	if err != nil {
		log.Printf("Not able to list functions: %s\n", out)
		log.Fatalln(err)
	}
	output := string(out[:])
	fmt.Println("This is done")
	return strings.TrimSuffix(output, "\n")
}

// fission env create --name nodejs --image fission/node-env
func (f FissionEnvironment) InitializeEnvironment() error {
	out := runFission(fmt.Sprintf("env create --name %s --image %s", f.Name, f.ImageName))
	fmt.Println(out)
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
	newFuncPath, created, downloadErr := f.downloadStub(httpMethod, o.ID)
	if downloadErr != nil {
		log.Fatalln(downloadErr)
		return downloadErr
	} else if created == false {
		// Function file with same name in same path already exists.
		return nil
	}
	funcCreate := fmt.Sprintf("function create --name %s --env %s --code %s", o.ID, f.Name, newFuncPath)
	out := runFission(funcCreate)
	fmt.Println(out)
	return nil
}

// DownloadStub downloads the FissionEnvironment's stub file but renamed to funcName
// httpDir is one of GET|POST|PUT|DELETE|HEAD which is a dir in the root of the project
func (f FissionEnvironment) downloadStub(httpMethod, funcName string) (string, bool, error) {
	functionFile := fmt.Sprintf("%s/%s%s", httpMethod, funcName, f.FileExtension)
	if _, err := os.Stat(functionFile); err == nil {
		return "", false, nil
	}
	curl := exec.Command("curl", "-Lo", functionFile, f.StubFileURL)
	_, err := curl.Output()
	if err != nil {
		fmt.Println(err)
		return "", false, errors.New(fmt.Sprintf("error downloading stub file %s", f.StubFileURL))
	}
	return functionFile, true, nil
}

// fission route create --method GET --url /hello --function hello
func (f FissionEnvironment) createRoute(httpMethod, funcName string, o *spec.Operation) error {
	routeCreate := fmt.Sprintf("route create --method %s --url %s --function %s", httpMethod, "/"+funcName, funcName)
	out := runFission(routeCreate)
	fmt.Println(out)
	return nil
}


func getFissionBinary() string {
	which := exec.Command("which", "fission")
	var path []byte
	path, err := which.Output()
	if err != nil {
		log.Println(`
		Fission tool is not installed. Assuming you're on MacOS, run below command
		curl -Lo fission https://github.com/fission/fission/releases/download/v0.4.0/fission-cli-osx
		&& chmod +x fission
		&& sudo mv fission /usr/local/bin/
		`)
		log.Fatalln(err)
	}
	pathStr := string(path[:])
	return strings.TrimSuffix(pathStr, "\n")
}
