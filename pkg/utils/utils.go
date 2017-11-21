package utils

import (
	"errors"
	"github.com/go-openapi/spec"
	"log"
	_ "fmt"
	"fmt"
	"os"
)

func Parse(p *spec.PathItemProps) (map[string]string, error) {
	_, err := getValidOperators(p)
	if err != nil {
		log.Fatalln(err)
	}
	//for _, pathItem := range paths {
	//	fmt.Println(spec.Operation.JSONLookup(pathItem, "Description"))
	//}
	return nil, nil
}

// Returns number of valid of
func getValidOperators(p *spec.PathItemProps) (map[string]spec.Operation, error) {
	validPaths := make(map[string]spec.Operation)
	if p.Get != nil {
		validPaths["get"] = *p.Get
	}
	if p.Post != nil {
		validPaths["post"] = *p.Post
	}
	if p.Put != nil {
		validPaths["put"] = *p.Put
	}
	if p.Patch != nil {
		fmt.Errorf("patch found but fission currently doesn't support that HTTP method so skipping")
	}
	if p.Delete != nil {
		validPaths["delete"] = *p.Delete
	}
	if p.Head != nil {
		validPaths["head"] = *p.Head
	}
	if p.Options != nil {
		validPaths["options"] = *p.Options
	}
	if len(validPaths) == 0 {
		return nil, errors.New("no valid operators for this path")
	}
	return validPaths, nil
}


func Cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("why can't we get current working dir?")
	}
	return dir
}