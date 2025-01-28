package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SimpleContract struct {
	contractapi.Contract
}

// Participant represents an individual in the system
type Participant struct {
	ParticipantID string `json:"participantID"`
	Role          string `json:"role"`  // e.g., "FirstResponder", "SecondInvestigator"
	MSPID         string `json:"mspID"` // e.g., "Org1MSP", "Org2MSP"
}

// Helper to check if caller is from the correct MSP
func (s *SimpleContract) checkMSP(ctx contractapi.TransactionContextInterface, allowedMSPs []string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to retrieve MSP ID: %v", err)
	}

	for _, allowedMSP := range allowedMSPs {
		if mspID == allowedMSP {
			return nil
		}
	}

	return fmt.Errorf("access denied: MSP ID '%s' is not authorized", mspID)
}

// RegisterParticipant allows adding a participant to the system
func (s *SimpleContract) RegisterParticipant(ctx contractapi.TransactionContextInterface, participantID, role string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to retrieve MSP ID: %v", err)
	}

	participant := Participant{
		ParticipantID: participantID,
		Role:          role,
		MSPID:         mspID,
	}

	participantJSON, err := json.Marshal(participant)
	if err != nil {
		return fmt.Errorf("failed to serialize participant: %v", err)
	}

	err = ctx.GetStub().PutState(participantID, participantJSON)
	if err != nil {
		return fmt.Errorf("failed to register participant: %v", err)
	}

	return nil
}

// GetParticipant retrieves a participant's details
func (s *SimpleContract) GetParticipant(ctx contractapi.TransactionContextInterface, participantID string) (*Participant, error) {
	participantJSON, err := ctx.GetStub().GetState(participantID)
	if err != nil || participantJSON == nil {
		return nil, fmt.Errorf("participant not found: %v", err)
	}

	var participant Participant
	err = json.Unmarshal(participantJSON, &participant)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize participant: %v", err)
	}

	return &participant, nil
}

// CreateEvidence restricted to Org1MSP (First Responder)
func (s *SimpleContract) CreateEvidence(ctx contractapi.TransactionContextInterface, evidenceID, description string) error {
	if err := s.checkMSP(ctx, []string{"Org1MSP"}); err != nil {
		return err
	}

	evidence := Evidence{
		EvidenceID:       evidenceID,
		Description:      description,
		Owner:            "FirstResponder",
		TransferHistory:  []string{"FirstResponder"},
		TimestampHistory: []string{time.Now().Format("2006-01-02 15:04:05")},
	}

	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, evidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to add evidence to ledger: %v", err)
	}

	caller, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to retrieve caller ID: %v", err)
	}
	fmt.Printf("Debug: Caller ID = %s\n", caller) // Debugging
	return s.logTransaction(ctx, "Create", evidenceID, caller)
}

// TransferEvidence allowed for both Org1MSP and Org2MSP
func (s *SimpleContract) TransferEvidence(ctx contractapi.TransactionContextInterface, evidenceID, newOwner string) error {
	if err := s.checkMSP(ctx, []string{"Org1MSP", "Org2MSP"}); err != nil {
		return err
	}

	evidenceJSON, err := ctx.GetStub().GetState(evidenceID)
	if err != nil || evidenceJSON == nil {
		return fmt.Errorf("evidence not found: %v", err)
	}

	var evidence Evidence
	err = json.Unmarshal(evidenceJSON, &evidence)
	if err != nil {
		return fmt.Errorf("failed to deserialize evidence: %v", err)
	}

	evidence.Owner = newOwner
	evidence.TransferHistory = append(evidence.TransferHistory, newOwner)
	evidence.TimestampHistory = append(evidence.TimestampHistory, time.Now().Format("2006-01-02 15:04:05"))

	updatedEvidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize updated evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, updatedEvidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to update evidence: %v", err)
	}

	caller, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to retrieve caller ID: %v", err)
	}
	fmt.Printf("Debug: Caller ID = %s\n", caller) // Debugging
	return s.logTransaction(ctx, "Transfer", evidenceID, caller)
}

// DeleteEvidence allowed for both Org1MSP and Org2MSP
func (s *SimpleContract) DeleteEvidence(ctx contractapi.TransactionContextInterface, evidenceID string) error {
	// Allow both Org1MSP and Org2MSP to delete evidence
	if err := s.checkMSP(ctx, []string{"Org1MSP", "Org2MSP"}); err != nil {
		return err
	}

	evidenceJSON, err := ctx.GetStub().GetState(evidenceID)
	if err != nil || evidenceJSON == nil {
		return fmt.Errorf("evidence not found: %v", err)
	}

	var evidence Evidence
	err = json.Unmarshal(evidenceJSON, &evidence)
	if err != nil {
		return fmt.Errorf("failed to deserialize evidence: %v", err)
	}

	evidence.Owner = "DELETED"
	updatedEvidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize updated evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, updatedEvidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to mark evidence as deleted: %v", err)
	}

	caller, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to retrieve caller ID: %v", err)
	}
	fmt.Printf("Debug: Caller ID = %s\n", caller) // Debugging
	return s.logTransaction(ctx, "Delete", evidenceID, caller)
}

// GetEvidenceDetails retrieves evidence details (accessible to all)
func (s *SimpleContract) GetEvidenceDetails(ctx contractapi.TransactionContextInterface, evidenceID string) (*Evidence, error) {
	evidenceJSON, err := ctx.GetStub().GetState(evidenceID)
	if err != nil || evidenceJSON == nil {
		return nil, fmt.Errorf("evidence not found: %v", err)
	}

	var evidence Evidence
	err = json.Unmarshal(evidenceJSON, &evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize evidence: %v", err)
	}

	return &evidence, nil
}

// Log transaction
func (s *SimpleContract) logTransaction(ctx contractapi.TransactionContextInterface, action, evidenceID, performedBy string) error {
	// Decode Base64-encoded Caller ID
	decodedID, err := base64.StdEncoding.DecodeString(performedBy)
	if err == nil {
		performedBy = string(decodedID)
		fmt.Printf("Debug: Decoded performedBy = '%s'\n", performedBy)
	} else {
		fmt.Printf("Debug: Failed to decode performedBy, using raw value = '%s'\n", performedBy)
	}

	// Lookup the participant in the ledger
	participant, err := s.GetParticipant(ctx, performedBy)
	role := "Unknown"
	if err == nil {
		role = participant.Role
	} else {
		mspID, _ := ctx.GetClientIdentity().GetMSPID()
		role = fmt.Sprintf("Unregistered (%s)", mspID)
		fmt.Printf("Debug: Participant not found for performedBy = '%s', using MSP ID = '%s'\n", performedBy, mspID)
	}

	// Log the transaction
	logID := fmt.Sprintf("LOG-%s", ctx.GetStub().GetTxID())
	log := TransactionLog{
		TransactionID: logID,
		Action:        action,
		EvidenceID:    evidenceID,
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		PerformedBy:   role,
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction log: %v", err)
	}

	return ctx.GetStub().PutState(logID, logJSON)
}

// GetTransactionLogs retrieves all transaction logs with serial numbers
func (s *SimpleContract) GetTransactionLogs(ctx contractapi.TransactionContextInterface) (string, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var logs []TransactionLog

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		// Filter for keys starting with "LOG-"
		if !strings.HasPrefix(queryResponse.Key, "LOG-") {
			continue
		}

		var log TransactionLog
		if err := json.Unmarshal(queryResponse.Value, &log); err == nil {
			logs = append(logs, log)
		}
	}

	// Sort logs by timestamp (ascending)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Timestamp < logs[j].Timestamp
	})

	// Build the formatted output with serial numbers
	var formattedLogs strings.Builder
	for i, log := range logs {
		logJSON, err := json.Marshal(log)
		if err != nil {
			return "", fmt.Errorf("failed to serialize log: %v", err)
		}
		formattedLogs.WriteString(fmt.Sprintf("%d - %s\n", i+1, logJSON))
	}

	return formattedLogs.String(), nil
}
