package model

import (
	"sync"
	"time"
)

//wait time if the queue is empty
var Polltime time.Duration
//maximum time for the given task to be completed since creation time
var MAX_PROCESS_TIME  = time.Minute*5

//state for task status
type State string

//various status of the task
const (
	UNTOUCHED State = "untouched"
	SUCCESS   State = "completed"
	FAILED    State = "failed"
	TIMEOUT   State = "timeout"
)

//model for taks object
type Task struct {
	Lock         sync.Mutex
	Id           string
	Status       State // untouched, completed, failed, timeout
	CreationTime time.Time // when was the task created
	TaskData     string // field containing data about the task
}
