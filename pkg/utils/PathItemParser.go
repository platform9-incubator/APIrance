package utils

import (
	"github.com/go-openapi/spec"
	"fmt"
)

func Parse(p *spec.PathItemProps) (map[string]string, error) {
	pathDescriptions := make(map[string]string)
	if p.Get != nil {
		pathDescriptions["GET"] = p.Get.Description
	}
	if p.Head != nil {
		pathDescriptions["HEAD"] = p.Get.Description
	}
	if p.Post != nil {
		pathDescriptions["POST"] = p.Get.Description
	}

	fmt.Println("ayyyy 😎")

	return pathDescriptions, nil
}
