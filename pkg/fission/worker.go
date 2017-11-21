package fission

import (
	"./env"
	"github.com/go-openapi/spec"
	"sync"
)

type Worker struct {
	WaitGroup   *sync.WaitGroup
	HttpMethod  string
	Operation   *spec.Operation
	Environment *env.FissionEnvironment
}

func Start(w *Worker) {
	defer w.WaitGroup.Done()
	w.Environment.Run(w.HttpMethod, w.Operation)
}
