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

// CreateEvidence adds new evidence
func (s *SimpleContract) CreateEvidence(ctx contractapi.TransactionContextInterface, evidenceID, description, owner string) error {
	evidence := Evidence{
		EvidenceID:       evidenceID,
		Description:      description,
		Owner:            owner,
		TransferHistory:  []string{owner},
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

	return s.logTransaction(ctx, "Create", evidenceID, owner)
}

// TransferEvidence updates ownership
func (s *SimpleContract) TransferEvidence(ctx contractapi.TransactionContextInterface, evidenceID, newOwner string) error {
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

// DeleteEvidence marks evidence as inactive
func (s *SimpleContract) DeleteEvidence(ctx contractapi.TransactionContextInterface, evidenceID string) error {
	evidenceJSON, err := ctx.GetStub().GetState(evidenceID)
	if err != nil || evidenceJSON == nil {
		return fmt.Errorf("evidence not found: %v", err)
	}

	var evidence Evidence
	err = json.Unmarshal(evidenceJSON, &evidence)
	if err != nil {
		return fmt.Errorf("failed to deserialize evidence: %v", err)
	}

	// Mark as deleted instead of removing
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

// GetEvidenceDetails retrieves evidence details
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

// logTransaction adds a transaction log
func (s *SimpleContract) logTransaction(ctx contractapi.TransactionContextInterface, action, evidenceID, performedBy string) error {
	logID := fmt.Sprintf("LOG-%s", ctx.GetStub().GetTxID()) // Use Fabric-provided TxID for consistency
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
