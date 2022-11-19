package worker

import "context"

type Worker interface {
	Name() string
	Run(context.Context) error
}
