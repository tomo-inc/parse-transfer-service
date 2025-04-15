package task

import (
	"context"
	"fmt"
)

type Task[T any] struct {
	resultChan chan T
	errorChan  chan error
}

func (t *Task[T]) Result() (result T, err error) {
	select {
	case result = <-t.resultChan:
	case err = <-t.errorChan:
	}

	return result, err
}

func NewAsyncTask[T any](ctx context.Context, tag string, fn func() (T, error)) *Task[T] {
	resultChan := make(chan T, 1)
	errorChan := make(chan error, 1)

	go asyncInner(ctx, tag, fn, resultChan, errorChan)

	return &Task[T]{
		resultChan: resultChan,
		errorChan:  errorChan,
	}

}

func asyncInner[T any](ctx context.Context, tag string, innerFn func() (T, error), resultChan chan T, errorChan chan error) {
	var (
		result T
		err    error
	)

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s panic: %v", tag, e)
		}

		if err != nil {
			errorChan <- err
		} else {
			resultChan <- result
		}
	}()

	result, err = innerFn()
}
