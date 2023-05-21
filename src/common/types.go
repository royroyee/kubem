package common

// Overview main
type Overview struct {
	NodeStatus NodeStatus `json:"node_status"`
	PodStatus  PodStatus  `json:"pod_status"`
}

type NodeStatus struct {
	NotReady []string `json:"not_ready"`
	Ready    []string `json:"ready"`
}

type PodStatus struct {
	Error   []string `json:"error"`
	Pending []string `json:"pending"`
	Running int      `json:"running"`
}

// Event
type Event struct {
	Created    string `json:"created"`
	EventLevel string `json:"event_level"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

type Count struct {
	Count int `json:"count"`
}

type NodeUsage struct {
	CpuUsage []int `json:"cpu_usage"`
	RamUsage []int `json:"ram_usage"`
}

type NodeOverview struct {
	Name     string  `json:"name"`
	CpuUsage float64 `json:"cpu_usage"`
	RamUsage float64 `json:"ram_usage"`
	IP       string  `json:"ip"`
	Status   string  `json:"status"`
}

type NodeInfo struct {
	OS                      string `json:"os"`
	HostName                string `json:"host_name"`
	IP                      string `json:"ip"`
	KubeletVersion          string `json:"kubelet_version"`
	ContainerRuntimeVersion string `json:"container_runtime_version"`
	NumContainers           int    `json:"num_containers"`
	CpuCores                int64  `json:"cpu_cores"`
	RamCapacity             int64  `json:"ram_capacity"`
	Status                  bool   `json:"status"`
}

type ControllerOverview struct {
	Namespace string   `json:"namespace"`
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Pods      []string `json:"pods"`
}

type ControllerDetail struct {
	TemplateContainers []string `json:"template_containers"`
	Volumes            []string `json:"volumes"`
}
