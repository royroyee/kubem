package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/royroyee/kubem/k8s"
	"net/http"

	"log"
)

// HTTP
type HTTPHandler struct {
	k8sHandler k8s.K8sHandler
}

func NewHTTPHandler(k8sHandler *k8s.K8sHandler) *HTTPHandler {
	httpHandler := &HTTPHandler{
		k8sHandler: *k8sHandler,
	}
	return httpHandler
}

func (httpHandler HTTPHandler) StartHTTPServer() {
	log.Println("Start HTTP Server .. ")

	r := httprouter.New()

	log.Println("Success to Start HTTP Server")

	// Overview
	r.GET("/overview/status", httpHandler.GetOverviewStatus)

	r.GET("/events", httpHandler.GetEvents) // Example : /events/?event=warning&page=1&per_page=10

	log.Fatal(http.ListenAndServe(":9000", r))

	log.Printf("Success to Start HTTP Server on port %d\n", 9000)

}
