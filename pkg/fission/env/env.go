package env

import (
	"github.com/go-openapi/spec"
	"os/exec"
	"fmt"
	"errors"
	"log"
	"os"
	"strings"
	_ "time"
)

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
		ImageName:     "fission/python-env",
		StubFileURL:   "https://raw.githubusercontent.com/fission/fission/master/examples/python/hello.py",
		FileExtension: ".py",
	},
}


// fission env create --name nodejs --image fission/node-env
func (f FissionEnvironment) InitializeEnvironment() error {
	envList := runFission("env list")
	if strings.Contains(envList, f.Name) {
		// The environment is already present
		log.Printf("Environment '%s' already present, skipping creation\n", f.Name)
		return nil
	}
	out := runFission(fmt.Sprintf("env create --name %s --image %s", f.Name, f.ImageName))
	fmt.Println(out)
	return nil
}

func (f FissionEnvironment) Run(httpMethod string, o *spec.Operation) error {
	f.createFunction(httpMethod, o)
	//time.Sleep(2000 * time.Millisecond) // 2 seconds gap to allow function to be available to create a route with
	// lowercase ID (name of function) since function names must be all lowercase but
	// routes can be created with a non-lowercase function name specified
	f.createRoute(httpMethod, strings.ToLower(o.ID))
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
		// TODO: add a flag to explicitly override function files
		return nil
	}
	funcCreate := fmt.Sprintf("function create --name %s --env %s --code %s",
		strings.ToLower(o.ID), f.Name, newFuncPath)
	out := runFission(funcCreate)
	fmt.Println(out)
	return nil
}

// DownloadStub downloads the FissionEnvironment's stub file but renamed to funcName
// httpDir is one of GET|POST|PUT|DELETE|HEAD which is a dir in the root of the project
func (f FissionEnvironment) downloadStub(httpMethod, funcName string) (string, bool, error) {
	// Relative path to where the stub file will be downloaded to, with its new name
	functionFile := fmt.Sprintf("fission/%s/%s%s", httpMethod, funcName, f.FileExtension)
	if _, err := os.Stat(functionFile); err == nil {
		fmt.Printf("[downloadStub] '%s' already exists, not downloading stub\n", functionFile)
		return "", false, nil
	}
	curl := exec.Command("curl", "-Lo", functionFile, f.StubFileURL)
	out, err := curl.CombinedOutput()
	if err != nil {
		fmt.Println(string(out[:]))
		return "", false, errors.New(fmt.Sprintf("error downloading stub file %s", f.StubFileURL))
	}
	return functionFile, true, nil
}

// fission route create --method GET --url /hello --function hello
func (f FissionEnvironment) createRoute(httpMethod, funcName string) error {
	routeCreate := fmt.Sprintf("route create --method %s --url %s --function %s", httpMethod, "/"+funcName, funcName)
	out := runFission(routeCreate)
	fmt.Println(out)
	return nil
}


// runFission takes a fission CLI command but without the 'fission', which is provided
// by 'getFissionBinary'. Returns results (if any) of command, and logs errors
func runFission(command string) string {
	args := fmt.Sprintf("%s %s", getFissionBinary(), command)
	fmt.Println("cmd: "+args)
	cmd := exec.Command("/usr/bin/bash", "-c", args)
	var out []byte
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error running fission command '%s': %s\n",args, string(out[:]))
		log.Fatalln(err)
	}
	output := string(out[:])
	return strings.TrimSuffix(output, "\n")
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
