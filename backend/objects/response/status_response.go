package response

type StatusResponse struct {
	Error   string `json:"error,omitempty"`
	Success string `json:"success,omitempty"`
}
