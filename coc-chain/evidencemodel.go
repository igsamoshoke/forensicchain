package main

type Evidence struct {
	EvidenceID       string   `json:"evidenceID"`
	Description      string   `json:"description"`
	Owner            string   `json:"owner"`
	TransferHistory  []string `json:"transferHistory"`
	TimestampHistory []string `json:"timestampHistory"`
}
