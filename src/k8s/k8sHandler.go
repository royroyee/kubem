package k8s

import (
	"context"
	cm "github.com/royroyee/kubem/common"
	"gopkg.in/mgo.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
)

// K8s
type K8sHandler struct {
	K8sClient       *kubernetes.Clientset
	MetricK8sClient *versioned.Clientset
	session         *mgo.Session
}

func NewK8sHandler() *K8sHandler {

	//In Cluster
	kh := &K8sHandler{

		K8sClient:       cm.ClientSetOutofCluster(),
		MetricK8sClient: cm.MetricClientSetOutofCluster(),
		session:         GetDBSession(),
	}

	////// Out of Cluster
	//kh := &K8sHandler{
	//	K8sClient:       cm.ClientSetOutofCluster(),
	//	MetricK8sClient: cm.MetricClientSetOutofCluster(),
	//	session:         GetDBSession(),
	//}

	return kh
}

func (kh K8sHandler) WatchEvents() {

	var result cm.Event

	watcher, err := kh.K8sClient.CoreV1().Events(metav1.NamespaceAll).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Failed to create event watcher: %v", err)
	}
	defer watcher.Stop()

	for watch := range watcher.ResultChan() {
		event, ok := watch.Object.(*v1.Event)
		if !ok {
			log.Println("Received non-Event object")
			continue
		}

		result.Created = event.LastTimestamp.Time.Format("2006-01-02 15:04")
		result.Name = event.InvolvedObject.Name
		result.Type = event.InvolvedObject.Kind
		result.Status = event.Reason
		result.Message = event.Message
		result.EventLevel = event.Type

		// TODO
		// kh.StoreEventInDB(result)
	}
}

// overview
func (kh K8sHandler) GetOverviewStatus() (cm.Overview, error) {
	var result cm.Overview

	ready, notReady, err := kh.nodeStatus()
	if err != nil {
		return result, err
	}
	running, pending, errorStatus, err := kh.podStatus()
	if err != nil {
		return result, err
	}

	result = cm.Overview{
		NodeStatus: cm.NodeStatus{
			NotReady: notReady,
			Ready:    ready,
		},
		PodStatus: cm.PodStatus{
			Error:   errorStatus,
			Pending: pending,
			Running: running,
		},
	}

	return result, nil
}
