package main

// Case represents a case containing related evidence and investigators
type Case struct {
	CaseID          string   `json:"caseID"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	EvidenceIDs     []string `json:"evidenceIDs"`
	InvestigatorIDs []string `json:"investigatorIDs"`
	Status          string   `json:"status"` // e.g., "Open", "Closed"
}
