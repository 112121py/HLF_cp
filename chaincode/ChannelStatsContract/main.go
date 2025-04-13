package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	chaincode, err := contractapi.NewChaincode(new(ChannelStatsContract))
	if err != nil {
		log.Panicf("Error creating ChannelStats chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting ChannelStats chaincode: %v", err)
	}
}
