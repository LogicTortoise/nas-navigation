package models

// SystemActionRequest represents a system action request
type SystemActionRequest struct {
	Action  string `json:"action" binding:"required"` // restart, shutdown, etc.
	Target  string `json:"target,omitempty"`          // optional target service name
	Confirm bool   `json:"confirm"`                   // confirmation flag
}

// SystemActionResult represents the result of a system action
type SystemActionResult struct {
	Success  bool   `json:"success"`
	Action   string `json:"action"`
	Message  string `json:"message"`
	Error    string `json:"error,omitempty"`
	Duration int64  `json:"duration_ms"`
}

// SystemInfo represents system information
type SystemInfo struct {
	Hostname    string `json:"hostname"`
	Platform    string `json:"platform"`
	Uptime      int64  `json:"uptime_seconds"`
	LoadAverage string `json:"load_average,omitempty"`
}

// ServiceRestartResult represents the result of restarting a single service
type ServiceRestartResult struct {
	ServiceID   uint   `json:"service_id"`
	ServiceName string `json:"service_name"`
	ScriptName  string `json:"script_name"`
	Success     bool   `json:"success"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
	Duration    int64  `json:"duration_ms"`
}

// BatchRestartResult represents the result of batch restarting all services
type BatchRestartResult struct {
	TotalCount   int                    `json:"total_count"`
	SuccessCount int                    `json:"success_count"`
	FailedCount  int                    `json:"failed_count"`
	Message      string                 `json:"message"`
	Results      []ServiceRestartResult `json:"results"`
}
