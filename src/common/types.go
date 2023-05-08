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
