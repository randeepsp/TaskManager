package executors

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"model"
	"queue"
	"time"
)

//starts executing the tasks from the queue
func StartExecutor(ctx context.Context, que *queue.Queue) {
	log.Printf("start executing tasks in queue")
	qr, _ := queue.NewQueueReader(que)
	//get the head of the queue
	head := qr.Next()
	for {
		if head == nil {
			//all items in queue have been processed once , wait and then start again
			//wait for a while else it will consume cpu unnecessarily
			time.Sleep(model.Polltime)
			head = qr.Next()
			continue
		}
		task, ok := head.Value.(*model.Task)
		if !ok {
			continue
		}
		executeTask(ctx, task)
		head = qr.Next()
	}
	log.Printf("completed exectuing tasks")
}

//executes the given task
func executeTask(ctx context.Context, task *model.Task) error {
	task.Lock.Lock()
	defer task.Lock.Unlock()

	//return if task has timedout or completed
	if (task.Status == model.SUCCESS) || (task.Status == model.TIMEOUT) {
		return nil
	}

	//mark the task as timed out if created more than max process time
	if time.Since(task.CreationTime) > model.MAX_PROCESS_TIME {
		task.Status = model.TIMEOUT
		return nil
	}

	//process only if task is untouched or failed
	//cleaner will remove the completed tasks
	if (task.Status == model.UNTOUCHED) || (task.Status == model.FAILED) {
		//attempt operation
		err := runOperationForTask()
		if err != nil {
			task.Status = model.FAILED
			log.Printf("operation failed for task %s", task.Id)
		} else {
			task.Status = model.SUCCESS
			log.Printf("operation succeeded for task %s", task.Id)
		}
		return nil
	}
	return nil
}

//simulates an operation ,failing and succeeding tasks randomly
func runOperationForTask() error {
	//trying an operation like sending an email
	randomNumber := rand.Intn(100)
	if randomNumber > 80 {
		//log.Printf("operation has failed")
		return  errors.New("operation failed")
	}
	return nil
}