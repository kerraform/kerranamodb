package dlock

import (
	"fmt"
	"strings"
)

const (
	deleminattor = "/"
)

type DLockID string

func From(table, key string) DLockID {
	return DLockID(fmt.Sprintf("%s%s%s", table, deleminattor, key))
}

func (id DLockID) String() string {
	return string(id)
}

func (id DLockID) Table() string {
	el := strings.Split(string(id), deleminattor)
	if len(el) < 1 {
		return ""
	}

	return el[0]
}

func (id DLockID) Key() string {
	el := strings.Split(string(id), deleminattor)
	if len(el) < 2 {
		return ""
	}

	return strings.Join(el[1:], deleminattor)
}
