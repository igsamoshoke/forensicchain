package main

type Participant struct {
	ParticipantID string `json:"participantID"`
	Role          string `json:"role"` // e.g., "First Responder", "Investigator"
}
