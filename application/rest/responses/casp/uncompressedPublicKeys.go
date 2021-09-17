package casp

type UncompressedPublicKeysResponse struct {
	TotalItems   int      `json:"totalItems"`
	Items        []string `json:"items"`
	Chains       []string `json:"chains"`
	AccountName  string   `json:"accountName"`
	AccountIndex int      `json:"accountIndex"`
}
