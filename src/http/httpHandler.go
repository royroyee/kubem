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

	r.GET("/overview/nodes/usage", httpHandler.GetNodeUsageOverview)

	// Event
	r.GET("/events", httpHandler.GetEvents) // Example : /events/?event=warning&page=1&per_page=10
	r.GET("/events/count", httpHandler.GetNumberOfEvents)

	// Nodes
	r.GET("/nodes", httpHandler.GetNodeOverview)
	r.GET("/node/usage/:name", httpHandler.GetNodeUsage)
	r.GET("/node/info/:name", httpHandler.GetNodeInfo)
	r.GET("/nodes/count", httpHandler.GetNumberOfNodes)

	// Workload
	r.GET("/workload/namespaces", httpHandler.GetNamespace)
	r.GET("/workload", httpHandler.GetControllersByFilter) // Filtering by Namespace, Type
	r.GET("/workload/count", httpHandler.GetNumberOfControllers)
	r.GET("/workload/info/:namespace/:name", httpHandler.GetControllerInfo)
	r.GET("/workload/conditions/:namespace/:name", httpHandler.GetConditions)
	r.GET("/workload/detail/:namespace/:name", httpHandler.GetControllerDetail)

	// Pod
	r.GET("/pod/info/:name", httpHandler.GetPodInfo) // Information of Pod (detail page)
	r.GET("/pod/usage/:name", httpHandler.GetPodUsage)
	r.GET("/pod/logs/:namespace/:name", httpHandler.GetLogsOfPod)

	log.Fatal(http.ListenAndServe(":9000", r))

	log.Printf("Success to Start HTTP Server on port %d\n", 9000)

}
