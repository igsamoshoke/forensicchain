package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(&SimpleContract{})
	if err != nil {
		panic("Error creating coc-chain contract: " + err.Error())
	}

	if err := chaincode.Start(); err != nil {
		panic("Error starting coc-chain contract: " + err.Error())
	}
}
