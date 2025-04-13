package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(FLTaskContract))
	if err != nil {
		log.Panicf("Error creating FLTask chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting FLTask chaincode: %v", err)
	}
}
