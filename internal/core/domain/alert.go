package domain

type AlertReceiver struct {
	Name string `json:"name"` // Name of the alert receiver
}

type AlertStatus struct {
	InhibitedBy []string `json:"inhibitedBy"` // List of inhibiting factors for the alert
	SilencedBy  []string `json:"silencedBy"`  // List of silence factors for the alert
	State       string   `json:"state"`       // Current state of the alert
}

type AlertLabels struct {
	AlertName string `json:"alertname"` // Name of the alert
}

type AlertAnnotations struct {
	Summary string `json:"summary"` // Summary of the alert
}

type Alert struct {
	Annotations  AlertAnnotations `json:"annotations"`  // Annotations for the alert
	EndsAt       string           `json:"endsAt"`       // Time at which the alert ends
	Fingerprint  string           `json:"fingerprint"`  // Unique identifier for the alert
	Receivers    []AlertReceiver  `json:"receivers"`    // List of receivers for the alert
	StartsAt     string           `json:"startsAt"`     // Time at which the alert starts
	Status       AlertStatus      `json:"status"`       // Status of the alert
	UpdatedAt    string           `json:"updatedAt"`    // Last update time of the alert
	GeneratorURL string           `json:"generatorURL"` // URL of the alert generator
	Labels       AlertLabels      `json:"labels"`       // Labels associated with the alert
}
