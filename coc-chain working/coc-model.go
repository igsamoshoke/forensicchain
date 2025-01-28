package main

type TransactionLog struct {
	TransactionID string `json:"transactionID"`
	Action        string `json:"action"`      // e.g., "Create", "Transfer", "Delete"
	EvidenceID    string `json:"evidenceID"`  // Links the transaction to specific evidence
	Timestamp     string `json:"timestamp"`   // ISO 8601 format preferred
	PerformedBy   string `json:"performedBy"` // ID of the participant
}
