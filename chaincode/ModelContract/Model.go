
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Model 結構
type Model struct {
	ID           string `json:"id"`
	Version      int    `json:"version"`
	Submitter    string `json:"submitter"`
	Params       string `json:"params"`
	ValidationStatus string `json:"validation_status"` // 驗證狀態 (Valid, Invalid, Pending)
	Timestamp    string `json:"timestamp"`
}

// ModelContract 定義
type ModelContract struct {
	contractapi.Contract
}

// SubmitModel 提交新模型
func (c *ModelContract) SubmitModel(ctx contractapi.TransactionContextInterface, id string, version int, submitter string, params string) error {
	model := Model{
		ID:           id,
		Version:      version,
		Submitter:    submitter,
		Params:       params,
		ValidationStatus: "Pending",  // 初始為 "Pending" 狀態
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	modelBytes, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal model: %v", err)
	}

	return ctx.GetStub().PutState(id, modelBytes)
}

// VerifyModel 驗證模型有效性
func (c *ModelContract) VerifyModel(ctx contractapi.TransactionContextInterface, id string, isValid bool) error {
	modelBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read model: %v", err)
	}
	if modelBytes == nil {
		return fmt.Errorf("model %s does not exist", id)
	}

	model := Model{}
	err = json.Unmarshal(modelBytes, &model)
	if err != nil {
		return fmt.Errorf("failed to unmarshal model: %v", err)
	}

	if model.ValidationStatus != "Pending" {
		return fmt.Errorf("model %s is already validated", id)
	}

	if isValid {
		model.ValidationStatus = "Valid"
	} else {
		model.ValidationStatus = "Invalid"
	}

	modelBytes, err = json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal updated model: %v", err)
	}

	return ctx.GetStub().PutState(id, modelBytes)
}

// RecordContribution 記錄模型貢獻者
func (c *ModelContract) RecordContribution(ctx contractapi.TransactionContextInterface, modelID string, contributor string) error {
	modelBytes, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return fmt.Errorf("failed to read model: %v", err)
	}
	if modelBytes == nil {
		return fmt.Errorf("model %s does not exist", modelID)
	}

	model := Model{}
	err = json.Unmarshal(modelBytes, &model)
	if err != nil {
		return fmt.Errorf("failed to unmarshal model: %v", err)
	}

	// Add contributor (simple implementation: append to params as a string list)
	model.Params += "," + contributor

	modelBytes, err = json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal updated model: %v", err)
	}

	return ctx.GetStub().PutState(modelID, modelBytes)
}

// GetModel 查詢模型資訊
func (c *ModelContract) GetModel(ctx contractapi.TransactionContextInterface, id string) (*Model, error) {
	modelBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read model: %v", err)
	}
	if modelBytes == nil {
		return nil, fmt.Errorf("model %s does not exist", id)
	}

	model := Model{}
	err = json.Unmarshal(modelBytes, &model)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal model: %v", err)
	}

	return &model, nil
}

// ListModels 列出所有模型
func (c *ModelContract) ListModels(ctx contractapi.TransactionContextInterface) ([]*Model, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get models: %v", err)
	}
	defer resultsIterator.Close()

	var models []*Model
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to read next model: %v", err)
		}

		model := Model{}
		err = json.Unmarshal(queryResponse.Value, &model)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal model: %v", err)
		}

		models = append(models, &model)
	}

	return models, nil
}
