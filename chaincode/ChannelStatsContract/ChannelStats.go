
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ChannelStats 定義全局資訊結構
type ChannelStats struct {
	TotalDataAmount     int    `json:"total_data_amount"`     // 總資料量
	TotalModelTrained   int    `json:"total_model_trained"`   // 總訓練模型次數
	TotalContributors   int    `json:"total_contributors"`    // 總貢獻者數量
	LastUpdated         string `json:"last_updated"`          // 最後更新時間
}

// ChannelStatsContract 定義合約
type ChannelStatsContract struct {
	contractapi.Contract
}

// InitializeChannelStats 初始化通道狀態
func (c *ChannelStatsContract) InitializeChannelStats(ctx contractapi.TransactionContextInterface) error {
	stats := ChannelStats{
		TotalDataAmount:   0,
		TotalModelTrained: 0,
		TotalContributors: 0,
		LastUpdated:       time.Now().Format(time.RFC3339),
	}

	statsBytes, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %v", err)
	}

	return ctx.GetStub().PutState("channel_stats", statsBytes)
}

// UpdateStats 更新通道統計資料
func (c *ChannelStatsContract) UpdateStats(ctx contractapi.TransactionContextInterface, dataAmount int, modelTrained bool, contributorAdded bool) error {
	statsBytes, err := ctx.GetStub().GetState("channel_stats")
	if err != nil {
		return fmt.Errorf("failed to read stats: %v", err)
	}
	if statsBytes == nil {
		return fmt.Errorf("channel stats do not exist")
	}

	stats := ChannelStats{}
	err = json.Unmarshal(statsBytes, &stats)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stats: %v", err)
	}

	// 更新資料量
	stats.TotalDataAmount += dataAmount

	// 訓練模型次數更新
	if modelTrained {
		stats.TotalModelTrained++
	}

	// 更新貢獻者數量
	if contributorAdded {
		stats.TotalContributors++
	}

	// 更新最後更新時間
	stats.LastUpdated = time.Now().Format(time.RFC3339)

	statsBytes, err = json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal updated stats: %v", err)
	}

	return ctx.GetStub().PutState("channel_stats", statsBytes)
}

// GetChannelStats 查詢通道統計資料
func (c *ChannelStatsContract) GetChannelStats(ctx contractapi.TransactionContextInterface) (*ChannelStats, error) {
	statsBytes, err := ctx.GetStub().GetState("channel_stats")
	if err != nil {
		return nil, fmt.Errorf("failed to read stats: %v", err)
	}
	if statsBytes == nil {
		return nil, fmt.Errorf("channel stats do not exist")
	}

	stats := ChannelStats{}
	err = json.Unmarshal(statsBytes, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal stats: %v", err)
	}

	return &stats, nil
}

// ListChannelStats 列出所有通道統計資料
func (c *ChannelStatsContract) ListChannelStats(ctx contractapi.TransactionContextInterface) ([]*ChannelStats, error) {
	// 此功能可以根據需求進行擴展，以下為簡單列出單一通道統計的範例
	stats, err := c.GetChannelStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel stats: %v", err)
	}

	return []*ChannelStats{stats}, nil
}

