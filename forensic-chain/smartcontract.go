package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing digital forensics operations
type SmartContract struct {
	contractapi.Contract
}

// CreateEvidence adds a new piece of evidence to the blockchain
func (s *SmartContract) CreateEvidence(ctx contractapi.TransactionContextInterface, evidenceID, creator, owner, description, caseID, fileHash string, fileSize int64, fileType string) error {
	role, exists, err := ctx.GetClientIdentity().GetAttributeValue("role")
	if err != nil {
		return errors.New("failed to retrieve participant role")
	}
	if !exists {
		return errors.New("participant role attribute does not exist")
	}
	if !CheckRole(role, "create") {
		return errors.New("participant is not authorized to create evidence")
	}

	existingEvidence, err := ctx.GetStub().GetState(evidenceID)
	if err != nil {
		return fmt.Errorf("failed to check if evidence exists: %v", err)
	}
	if existingEvidence != nil {
		return errors.New("evidence with this ID already exists")
	}

	evidence := Evidence{
		EvidenceID:    evidenceID,
		Creator:       creator,
		Owner:         owner,
		Description:   description,
		CaseID:        caseID,
		TransferChain: []string{creator},
		TransferTime:  []string{},
		FileHash:      fileHash,
		FileSize:      fileSize,
		FileType:      fileType,
	}

	evidenceBytes, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to marshal evidence: %v", err)
	}

	err = ctx.GetStub().PutState(evidenceID, evidenceBytes)
	if err != nil {
		return fmt.Errorf("failed to store evidence: %v", err)
	}

	timestamp := time.Now().Format(time.RFC3339)
	return s.LogAction(ctx, evidenceID, creator, "Created", timestamp)
}

// LogAction logs an action for evidence auditing
func (s *SmartContract) LogAction(ctx contractapi.TransactionContextInterface, evidenceID, participantID, actionType, timestamp string) error {
	log := ActionLog{
		ActionID:      ctx.GetStub().GetTxID(),
		EvidenceID:    evidenceID,
		ParticipantID: participantID,
		ActionType:    actionType,
		Timestamp:     timestamp,
	}

	logBytes, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal action log: %v", err)
	}

	return ctx.GetStub().PutState(log.ActionID, logBytes)
}

// Main function
func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("Error creating smart contract: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error starting smart contract: %v", err))
	}
}
