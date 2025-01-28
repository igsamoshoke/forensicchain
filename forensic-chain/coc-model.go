package main

// ActionLog records every action for auditing purposes
type ActionLog struct {
	ActionID      string `json:"actionID"`
	EvidenceID    string `json:"evidenceID"`
	ParticipantID string `json:"participantID"`
	ActionType    string `json:"actionType"` // e.g., "Created", "Transferred", "Deleted"
	Timestamp     string `json:"timestamp"`
}
