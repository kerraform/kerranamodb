package api

import (
	"fmt"

	"github.com/kerraform/kerranamodb/internal/id"
)

type GetInput struct {
	TableName           string
	Key                 map[string]map[string]string
	ConditionExpression string
}

func (i *GetInput) GetInfo() (string, error) {
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

func (i *GetInput) GetLockID() (id.LockID, error) {
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
