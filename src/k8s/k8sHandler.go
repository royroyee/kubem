package k8s

import (
	"bufio"
	"context"
	"fmt"
	cm "github.com/royroyee/kubem/common"
	"gopkg.in/mgo.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"regexp"
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

func (kh K8sHandler) GetNodeInfo(nodeName string) (cm.NodeInfo, error) {

	var result cm.NodeInfo

	node, err := kh.K8sClient.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return result, err
	}

	result.OS = node.Status.NodeInfo.OSImage
	result.HostName = node.ObjectMeta.Name
	result.IP = node.Status.Addresses[0].Address
	result.Status = isNodeReady(node)

	result.KubeletVersion = node.Status.NodeInfo.KubeletVersion
	result.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion

	pods, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{FieldSelector: "spec.nodeName=" + result.HostName})
	if err != nil {
		log.Println(err)
		return result, err
	}
	numContainers := 0
	for _, pod := range pods.Items {
		numContainers += len(pod.Spec.Containers)
	}
	result.NumContainers = numContainers
	capacity := node.Status.Capacity
	result.CpuCores = capacity.Cpu().Value()
	result.RamCapacity = node.Status.Capacity.Memory().Value() / 1024 / 1024 / 1024

	return result, nil

}

func (kh K8sHandler) NumberOfNodes() (cm.Count, error) {
	var result cm.Count

	nodes, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return result, err
	}

	result.Count = len(nodes.Items)
	return result, err
}

func (kh K8sHandler) GetControllerDetail(namespace string, name string) (cm.ControllerDetail, error) {

	result, err := kh.GetVolumesOfController(namespace, name)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (kh K8sHandler) GetControllerInfo(controllerType string, namespace string, controllerName string) (cm.ControllerInfo, error) {
	var result cm.ControllerInfo
	var limits, volumes, mounts, envs, labels []string
	var controlleredByName string

	if controllerType == "deployment" {
		controller, err := kh.K8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}

		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, fmt.Sprintf("%s:%s", env.Name, env.Value))
			}

			for key, value := range controller.Labels {
				labels = append(labels, fmt.Sprintf("%s=%s", key, value))
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "daemonset" {
		controller, err := kh.K8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}

		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "staefulset" {
		controller, err := kh.K8sClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "job" {
		controller, err := kh.K8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "cronjob" {
		controller, err := kh.K8sClient.BatchV1().CronJobs(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.JobTemplate.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.JobTemplate.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else if controllerType == "replicaset" {
		controller, err := kh.K8sClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), controllerName, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		container := controller.Spec.Template.Spec.Containers

		if len(container) > 0 {
			for _, container := range controller.Spec.Template.Spec.Containers {
				for resourceName, resourceLimit := range container.Resources.Limits {
					limits = append(limits, fmt.Sprintf("%s=%s", resourceName, resourceLimit.String()))
				}
			}
			for _, volume := range controller.Spec.Template.Spec.Volumes {
				volumes = append(volumes, volume.Name)
			}
			for _, volumeMount := range controller.Spec.Template.Spec.Containers[0].VolumeMounts {
				mounts = append(mounts, volumeMount.Name)
			}
			for _, env := range controller.Spec.Template.Spec.Containers[0].Env {
				envs = append(envs, env.Name)
			}
			if len(controller.OwnerReferences) > 0 {
				controlleredBy := controller.OwnerReferences[0]
				controlleredByName = controlleredBy.Name
			}

		}

	} else {
		err := fmt.Errorf("Invalid Controller Type %v", controllerType)
		return result, err
	}

	result.Limits = limits
	result.Environment = envs
	result.Mounts = mounts
	result.Volumes = volumes
	result.Labels = labels
	result.ControlledBy = controlleredByName

	return result, nil
}

func (kh K8sHandler) GetConditions(controllerType string, namespace string, name string) ([]cm.Conditions, error) {

	var result []cm.Conditions

	if controllerType == "deployment" {

		deployment, err := kh.K8sClient.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		conditions := deployment.Status.Conditions
		if len(conditions) > 0 {
			for _, condition := range conditions {
				//result.Type = append(result.Type, string(condition.Type))
				//result.Status = append(result.Status, string(condition.Status))
				//result.Reason = append(result.Reason, condition.Reason)

				result = append(result, cm.Conditions{
					Type:   string(condition.Type),
					Status: string(condition.Status),
					Reason: (condition.Reason),
				})
			}
		}
	} else if controllerType == "daemonset" {
		daemonset, err := kh.K8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		conditions := daemonset.Status.Conditions
		if len(conditions) > 0 {
			for _, condition := range conditions {
				result = append(result, cm.Conditions{
					Type:   string(condition.Type),
					Status: string(condition.Status),
					Reason: (condition.Reason),
				})
			}
		}

	} else if controllerType == "staefulset" {
		statefulset, err := kh.K8sClient.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		conditions := statefulset.Status.Conditions
		if len(conditions) > 0 {
			for _, condition := range conditions {
				result = append(result, cm.Conditions{
					Type:   string(condition.Type),
					Status: string(condition.Status),
					Reason: (condition.Reason),
				})
			}
		}

	} else if controllerType == "job" {
		job, err := kh.K8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		conditions := job.Status.Conditions
		if len(conditions) > 0 {
			for _, condition := range conditions {
				result = append(result, cm.Conditions{
					Type:   string(condition.Type),
					Status: string(condition.Status),
					Reason: (condition.Reason),
				})
			}
		}

	} else if controllerType == "replicaset" {
		replicaset, err := kh.K8sClient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return result, err
		}
		conditions := replicaset.Status.Conditions
		if len(conditions) > 0 {
			for _, condition := range conditions {
				result = append(result, cm.Conditions{
					Type:   string(condition.Type),
					Status: string(condition.Status),
					Reason: (condition.Reason),
				})
			}
		}

	} else {
		err := fmt.Errorf("Invalid Controller Type %v", controllerType)
		return result, err
	}
	return result, nil
}

func (kh K8sHandler) GetLogsOfPod(namespace string, podName string) ([]string, error) {
	var result []string

	// create options for retrieving the logs
	options := &v1.PodLogOptions{
		Timestamps: true,
		TailLines:  new(int64),
	}
	*options.TailLines = 30

	// get the logs for the specified pod
	req := kh.K8sClient.CoreV1().Pods(namespace).GetLogs(podName, options)
	logs, err := req.Stream(context.Background())
	if err != nil {
		return result, err
	}
	defer logs.Close()

	// read the logs and format them with timestamps and pod name
	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		line := scanner.Text()

		// extract the "$date" field from the JSON object in the log line
		re := regexp.MustCompile(`\{"\$date":"([^"]+)"\}`)
		match := re.FindStringSubmatch(line)
		var dateStr string
		if len(match) == 2 {
			dateStr = match[1]
		}

		// format the log line with the timestamp and pod name
		formatted := fmt.Sprintf("%s [%s] %s", dateStr, podName, line)
		result = append(result, formatted)
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	// return the result slice
	return result, nil
}
