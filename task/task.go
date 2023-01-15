package task

import (
	pool "github.com/ennesuysal/go-thread-pooling"
)

func (t *Task) Execute(param interface{}) error {
	return t.executeFunc(param)
}

func (t *Task) OnFailure(e error) {
	t.failureFunc(e)
}

type Task struct {
	executeFunc func(interface{}) error
	failureFunc func(error)
}

func NewTask(executeFunc func(interface{}) error, failureFunc func(error), parameters interface{}) pool.Task {
	return pool.Task{
		Exec: &Task{
			executeFunc: executeFunc,
			failureFunc: failureFunc,
		},
		Parameters: parameters,
	}
}
