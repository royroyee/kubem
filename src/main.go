package src

import (
	"github.com/royroyee/kubem/http"
	"github.com/royroyee/kubem/k8s"
	"log"
	"sync"
)

var handlers Handlers

type Handlers struct {
	k8sHandler  *k8s.K8sHandler
	httpHandler *http.HTTPHandler
}

func initHandlers() {
	handlers.k8sHandler = k8s.NewK8sHandler()
	handlers.httpHandler = http.NewHTTPHandler(handlers.k8sHandler)
}

func main() {

	log.Println("Welcome to kubem!")

	// Handlers
	initHandlers()

	var wg sync.WaitGroup
	wg.Add(3)

	// Start DB Session
	go handlers.k8sHandler.DBSession()

	// Start HTTP Servers
	go handlers.httpHandler.StartHTTPServer()

	go handlers.k8sHandler.WatchEvents()

	wg.Wait()
	log.Println("kubem  finished. Bye.")
}
