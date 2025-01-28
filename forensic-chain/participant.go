package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Participant represents a forensic investigator or an entity interacting with evidence
type Participant struct {
	ParticipantID string `json:"participantID"`
	Name          string `json:"name"`
	Role          string `json:"role"` // e.g., "first responder", "second investigator", etc.
}

// AddParticipant adds a new participant to the blockchain
func (s *SmartContract) AddParticipant(ctx contractapi.TransactionContextInterface, participantID, name, role string) error {
	existingParticipant, err := ctx.GetStub().GetState(participantID)
	if err != nil {
		return fmt.Errorf("failed to check if participant exists: %v", err)
	}
	if existingParticipant != nil {
		return errors.New("participant with this ID already exists")
	}

	participant := Participant{
		ParticipantID: participantID,
		Name:          name,
		Role:          role,
	}

	participantBytes, err := json.Marshal(participant)
	if err != nil {
		return fmt.Errorf("failed to marshal participant: %v", err)
	}

	return ctx.GetStub().PutState(participantID, participantBytes)
}
