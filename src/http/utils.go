package http

import (
	"github.com/julienschmidt/httprouter"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"strconv"
)

func (httpHandler HTTPHandler) GetOverviewStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	overview, err := httpHandler.k8sHandler.GetOverviewStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&overview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetEvents(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	eventType := r.URL.Query().Get("event")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	// Get the data from db
	events, err := httpHandler.k8sHandler.GetEvents(eventType, page, perPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNumberOfEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	eventLevel := r.URL.Query().Get("level")
	count, err := httpHandler.k8sHandler.NumberOfEvents(eventLevel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeUsageOverview(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	nodeUsage, err := httpHandler.k8sHandler.GetNodeUsageAvg()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeUsage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeOverview(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	nodeOverview, err := httpHandler.k8sHandler.GetNodeOverview(page, perPage)
	result, err := json.Marshal(&nodeOverview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeUsage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	nodeUsage, err := httpHandler.k8sHandler.GetNodeUsage(ps.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeUsage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNodeInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	nodeInfo, err := httpHandler.k8sHandler.GetNodeInfo(params.ByName("name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&nodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNumberOfNodes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	count, err := httpHandler.k8sHandler.NumberOfNodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNamespace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	namespaceList, err := httpHandler.k8sHandler.GetNamespaceName()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&namespaceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetControllersByFilter(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	namespace := r.URL.Query().Get("namespace")
	controller := r.URL.Query().Get("controller")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		perPage = 10
	}

	controllers, err := httpHandler.k8sHandler.GetControllersByFilter(namespace, controller, page, perPage)
	result, err := json.Marshal(&controllers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (httpHandler HTTPHandler) GetNumberOfControllers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	namespace := r.URL.Query().Get("namespace")
	controllerType := r.URL.Query().Get("type")
	count, err := httpHandler.k8sHandler.NumberOfControllers(namespace, controllerType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := json.Marshal(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
