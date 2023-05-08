package k8s

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func (kh K8sHandler) nodeStatus() (ready []string, notReady []string, err error) {

	nodeList, err := kh.K8sClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return ready, notReady, err
	}

	for _, node := range nodeList.Items {
		if isNodeReady(&node) {
			ready = append(ready, node.GetName())
		} else {
			notReady = append(ready, node.GetName())
		}
	}
	return ready, notReady, nil
}

func (kh K8sHandler) podStatus() (running int, pending []string, errorStatus []string, err error) {

	running = 0

	podList, err := kh.K8sClient.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return running, pending, errorStatus, err
	}

	for _, pod := range podList.Items {
		switch pod.Status.Phase {
		case corev1.PodPending:
			pending = append(pending, pod.Name)
		case corev1.PodRunning:
			running++
		case corev1.PodSucceeded:
			running++
		case corev1.PodFailed:
			errorStatus = append(errorStatus, pod.Name)
		default:
			errorStatus = append(errorStatus, pod.Name)
		}
	}
	return running, pending, errorStatus, nil
}

func isNodeReady(node *corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func NodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			return "Ready"
		}
	}
	return "Not Ready"
}
