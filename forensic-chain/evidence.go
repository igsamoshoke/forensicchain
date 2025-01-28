package main

// Evidence represents the structure of digital evidence on the blockchain
type Evidence struct {
	EvidenceID    string   `json:"evidenceID"`
	Creator       string   `json:"creator"`
	Owner         string   `json:"owner"`
	Description   string   `json:"description"`
	CaseID        string   `json:"caseID"`
	TransferChain []string `json:"transferChain"`
	TransferTime  []string `json:"transferTime"`
	FileHash      string   `json:"fileHash"`
	FileSize      int64    `json:"fileSize"`
	FileType      string   `json:"fileType"`
}
