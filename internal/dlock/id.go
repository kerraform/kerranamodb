package dlock

import (
	"fmt"
	"strings"
)

type DLockID string

func From(table, key string) DLockID {
	return DLockID(fmt.Sprintf("%s/%s", table, key))
}

func (id DLockID) String() string {
	return string(id)
}

func (id DLockID) Table() string {
	el := strings.Split(string(id), "/")
	if len(el) < 1 {
		return ""
	}

	return el[0]
}

func (id DLockID) Key() string {
	el := strings.Split(string(id), "/")
	if len(el) < 2 {
		return ""
	}

	return el[1]
}
