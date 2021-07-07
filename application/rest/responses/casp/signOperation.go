package casp

import "time"

type SignOperationResponse struct {
	OperationID string    `json:"operationID"`
	Kind        string    `json:"kind"`
	Status      string    `json:"status"`
	StatusText  string    `json:"statusText"`
	CreatedAt   time.Time `json:"createdAt"`
	VaultID     string    `json:"vaultID"`
	Description string    `json:"description,omitempty"`
	IsApproved  bool      `json:"isApproved"`
	AccountID   string    `json:"accountID"`
	Groups      []struct {
		Name    string `json:"name"`
		Members []struct {
			ApprovedAt           string `json:"approvedAt"`
			Id                   string `json:"id"`
			IsApproved           bool   `json:"isApproved"`
			Name                 string `json:"name"`
			Status               string `json:"status"`
			ApprovalGroupAccount struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"approvalGroupAccount"`
			Label        string `json:"label"`
			IsGlobal     bool   `json:"isGlobal"`
			IsActive     bool   `json:"isActive"`
			GlobalStatus string `json:"globalStatus"`
			IsOffline    bool   `json:"isOffline"`
		} `json:"members"`
		RequiredApprovals int  `json:"requiredApprovals"`
		Order             int  `json:"order"`
		DeactivateAllowed bool `json:"deactivateAllowed"`
		IsOffline         bool `json:"isOffline"`
	} `json:"groups"`
	VaultName           string   `json:"vaultName"`
	PublicKeys          []string `json:"publicKeys"`
	DataToSign          []string `json:"dataToSign"`
	Signatures          []string `json:"signatures,omitempty"`
	V                   []int    `json:"v,omitempty"`
	LedgerHashAlgorithm string   `json:"ledgerHashAlgorithm"`
	CollectedData       struct {
		CollectionComplete   bool          `json:"collectionComplete"`
		DataCollectionGroups []interface{} `json:"dataCollectionGroups"`
	} `json:"collectedData"`
}
