package main

type TransactionLog struct {
	TransactionID string `json:"transactionID"`
	Action        string `json:"action"`
	EvidenceID    string `json:"evidenceID"`
	Timestamp     string `json:"timestamp"`
	PerformedBy   string `json:"performedBy"`
}
