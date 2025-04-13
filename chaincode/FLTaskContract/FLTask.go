


package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Task 結構
type Task struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`
	Round     int      `json:"round"`
	Participants []string `json:"participants"`
	Params    string   `json:"params"`
	Timestamp string   `json:"timestamp"`
}

// FLTaskContract 定義
type FLTaskContract struct {
	contractapi.Contract
}

// CreateTask 創建新任務
func (c *FLTaskContract) CreateTask(ctx contractapi.TransactionContextInterface, id string, taskType string, round int, participants []string, params string) error {
	task := Task{
		ID:           id,
		Type:         taskType,
		Round:        round,
		Participants: participants,
		Params:       params,
		Timestamp:    time.Now().Format(time.RFC3339),
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	return ctx.GetStub().PutState(id, taskBytes)
}

// UpdateTask 更新任務資訊
func (c *FLTaskContract) UpdateTask(ctx contractapi.TransactionContextInterface, id string, taskType string, round int, participants []string, params string) error {
	taskBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read task: %v", err)
	}
	if taskBytes == nil {
		return fmt.Errorf("task %s does not exist", id)
	}

	task := Task{}
	err = json.Unmarshal(taskBytes, &task)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task: %v", err)
	}

	task.Type = taskType
	task.Round = round
	task.Participants = participants
	task.Params = params
	task.Timestamp = time.Now().Format(time.RFC3339)

	taskBytes, err = json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal updated task: %v", err)
	}

	return ctx.GetStub().PutState(id, taskBytes)
}

// GetTask 查詢任務資訊
func (c *FLTaskContract) GetTask(ctx contractapi.TransactionContextInterface, id string) (*Task, error) {
	taskBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read task: %v", err)
	}
	if taskBytes == nil {
		return nil, fmt.Errorf("task %s does not exist", id)
	}

	task := Task{}
	err = json.Unmarshal(taskBytes, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %v", err)
	}

	return &task, nil
}

// ListTasks 列出所有任務
func (c *FLTaskContract) ListTasks(ctx contractapi.TransactionContextInterface) ([]*Task, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer resultsIterator.Close()

	var tasks []*Task
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to read next task: %v", err)
		}

		task := Task{}
		err = json.Unmarshal(queryResponse.Value, &task)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal task: %v", err)
		}

		tasks = append(tasks, &task)
	}

	return tasks, nil
}

