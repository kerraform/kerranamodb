package api

import (
	"fmt"

	"github.com/kerraform/kerranamodb/internal/id"
)

const (
	InfoKey   = "Info"
	LockIDKey = "LockID"
	SKey      = "S"
)

type PutInput struct {
	TableName           string `json:"TableName"`
	Item                map[string]map[string]string
	ConditionExpression string
}

func (i *PutInput) GetInfo() (string, error) {
	k, ok := i.Item[InfoKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", LockIDKey)
	}

	res, ok := k[SKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", SKey)
	}

	return res, nil
}

func (i *PutInput) GetLockID() (id.LockID, error) {
	k, ok := i.Item[LockIDKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", LockIDKey)
	}

	res, ok := k[SKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", SKey)
	}

	return id.LockID(res), nil
}
