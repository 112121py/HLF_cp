
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ModelVerification 定義模型驗證資料結構
type ModelVerification struct {
	ModelID        string `json:"model_id"`
	ModelHash      string `json:"model_hash"`
	Signer         string `json:"signer"`
	Signature      string `json:"signature"`
	VerificationStatus string `json:"verification_status"` // 驗證狀態 (Verified/Failed)
	TrustLevel     string `json:"trust_level"`          // 可信度 (High/Medium/Low)
	Timestamp      string `json:"timestamp"`            // 提交時間
}

// ZeroTrustEndorseContract 定義背書合約
type ZeroTrustEndorseContract struct {
	contractapi.Contract
}

// EndorseModel 背書模型
func (c *ZeroTrustEndorseContract) EndorseModel(ctx contractapi.TransactionContextInterface, modelID string, modelHash string, signer string, signature string) error {
	// 計算模型哈希值
	calculatedHash := sha256.Sum256([]byte(modelHash))

	// 模型驗證資料結構
	verification := ModelVerification{
		ModelID:        modelID,
		ModelHash:      fmt.Sprintf("%x", calculatedHash),
		Signer:         signer,
		Signature:      signature,
		VerificationStatus: "Verified",  // 預設狀態為 Verified
		TrustLevel:     "High",         // 根據簽名者的可信度設定
		Timestamp:      time.Now().Format(time.RFC3339),
	}

	// 驗證簽名（此處可擴展實現具體的簽名驗證邏輯）
	// 模擬簽名驗證，實際可根據需求使用數位證書或其他方式
	if signer == "" || signature == "" {
		verification.VerificationStatus = "Failed"
		verification.TrustLevel = "Low"
	}

	// 將背書信息保存至區塊鏈
	verificationBytes, err := json.Marshal(verification)
	if err != nil {
		return fmt.Errorf("failed to marshal verification data: %v", err)
	}

	return ctx.GetStub().PutState(modelID, verificationBytes)
}

// GetModelVerification 查詢模型的背書信息
func (c *ZeroTrustEndorseContract) GetModelVerification(ctx contractapi.TransactionContextInterface, modelID string) (*ModelVerification, error) {
	// 從區塊鏈查詢模型背書資訊
	verificationBytes, err := ctx.GetStub().GetState(modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get model verification data: %v", err)
	}
	if verificationBytes == nil {
		return nil, fmt.Errorf("model verification data not found for modelID %s", modelID)
	}

	// 反序列化模型驗證資料
	verification := ModelVerification{}
	err = json.Unmarshal(verificationBytes, &verification)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal verification data: %v", err)
	}

	return &verification, nil
}

// ListAllVerifications 列出所有模型的背書資料
func (c *ZeroTrustEndorseContract) ListAllVerifications(ctx contractapi.TransactionContextInterface) ([]*ModelVerification, error) {
	// 此功能可根據需求進行擴展，以下為簡單列出單一模型背書的範例
	// 也可以增加遍歷區塊鏈中所有模型的功能
	var allVerifications []*ModelVerification

	// 這裡假設只有一個模型背書的資料
	verification, err := c.GetModelVerification(ctx, "example_model_id")
	if err != nil {
		return nil, fmt.Errorf("failed to get model verification: %v", err)
	}

	allVerifications = append(allVerifications, verification)
	return allVerifications, nil
}