package casp

type UncompressedPublicKeysResponse struct {
	PublicKeys   []string `json:"publicKeys"`
	Chains       []string `json:"chains"`
	AccountName  string   `json:"accountName"`
	AccountIndex int      `json:"accountIndex"`
}
