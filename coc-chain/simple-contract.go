package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SimpleContract struct {
	contractapi.Contract
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

// CreateEvidence restricted to Org1MSP (First Responder)
func (s *SimpleContract) CreateEvidence(ctx contractapi.TransactionContextInterface, evidenceID, description string) error {
	// Ensure this function is only accessible to Org1MSP (First Responder)
	if err := s.checkMSP(ctx, []string{"Org1MSP"}); err != nil {
		return err
	}

	evidence := Evidence{
		EvidenceID:       evidenceID,
		Description:      description,
		Owner:            "FirstResponder",
		TransferHistory:  []string{"FirstResponder"},
		TimestampHistory: []string{time.Now().Format(time.RFC3339)},
	}

	evidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, evidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to add evidence to ledger: %v", err)
	}

	return s.logTransaction(ctx, "Create", evidenceID, "FirstResponder")
}

// TransferEvidence allowed for both Org1MSP and Org2MSP
func (s *SimpleContract) TransferEvidence(ctx contractapi.TransactionContextInterface, evidenceID, newOwner string) error {
	// Restrict to Org1MSP and Org2MSP
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
	evidence.TimestampHistory = append(evidence.TimestampHistory, time.Now().Format(time.RFC3339))

	updatedEvidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize updated evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, updatedEvidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to update evidence: %v", err)
	}

	return s.logTransaction(ctx, "Transfer", evidenceID, newOwner)
}

// DeleteEvidence allowed for both Org1MSP and Org2MSP
func (s *SimpleContract) DeleteEvidence(ctx contractapi.TransactionContextInterface, evidenceID string) error {
	// Restrict to Org2MSP only (Second Responder)
	if err := s.checkMSP(ctx, []string{"Org2MSP"}); err != nil {
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

	// Mark as deleted
	evidence.Owner = "DELETED"
	updatedEvidenceJSON, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to serialize updated evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, updatedEvidenceJSON)
	if err != nil {
		return fmt.Errorf("failed to mark evidence as deleted: %v", err)
	}

	caller, _ := ctx.GetClientIdentity().GetID()
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
	logID := fmt.Sprintf("LOG-%s", ctx.GetStub().GetTxID())
	log := TransactionLog{
		TransactionID: logID,
		Action:        action,
		EvidenceID:    evidenceID,
		Timestamp:     time.Now().Format(time.RFC3339),
		PerformedBy:   performedBy,
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to serialize transaction log: %v", err)
	}

	return ctx.GetStub().PutState(logID, logJSON)
}

// GetTransactionLogs retrieves all transaction logs
func (s *SimpleContract) GetTransactionLogs(ctx contractapi.TransactionContextInterface) ([]TransactionLog, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var logs []TransactionLog
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	return logs, nil
}
