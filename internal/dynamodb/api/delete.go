package api

import (
	"fmt"

	"github.com/kerraform/kerranamodb/internal/id"
)

type DeleteInput struct {
	TableName string `json:"TableName"`
	Key       map[string]map[string]string
}

func (i *DeleteInput) GetInfo() (string, error) {
	k, ok := i.Key[InfoKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", LockIDKey)
	}

	res, ok := k[SKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", SKey)
	}

	return res, nil
}

func (i *DeleteInput) GetLockID() (id.LockID, error) {
	k, ok := i.Key[LockIDKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", LockIDKey)
	}

	res, ok := k[SKey]
	if !ok {
		return "", fmt.Errorf("%s not exist", SKey)
	}

	return id.LockID(res), nil
}
