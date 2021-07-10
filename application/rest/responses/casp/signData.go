package casp

type PostSignDataResponse struct {
	OperationID string `json:"operationID"`
}

type ErrorResponse struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
}
