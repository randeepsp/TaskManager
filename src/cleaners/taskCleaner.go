package cleaners

import (
	"context"
	"log"
	"model"
	"queue"
	"sync"
	"time"
)

//starts cleaning the queue for timedout or completed tasks
func StartCleaner(ctx context.Context, wg *sync.WaitGroup, queue *queue.Queue){
	defer wg.Done()
	log.Printf("begin checking queue for cleaning up of tasks")
	cleanUpTasks(ctx, queue)
}

//cleans up the given task
func cleanUpTasks(ctx context.Context, que *queue.Queue){
	log.Printf("called cleanup tasks on queue")
	qr, _ := queue.NewQueueReader(que)
	//get the head of the queue
	head := qr.Next()
	for {
		if head == nil {
			//all items in queue have been processed once , wait and then start again
			//dont need to run as it will consume cpu unnecessarily
			time.Sleep(model.Polltime)
			head = qr.Next()
			continue
		}
		task, ok := head.Value.(*model.Task)
		if !ok {
			continue
		}
		//check if task needs to be removed or added to end of queue
		pushBack, remove := checkTaskStatus(ctx, task)
		if pushBack{
			que.RemoveElement(head)
			que.Enqueue(head)
		}
		if remove {
			que.RemoveElement(head)
		}
		head = qr.Next()
	}
}

//check the status of the task, either we remove or push it at the back
func checkTaskStatus(ctx context.Context, task *model.Task) (pushBack bool, remove bool){
	task.Lock.Lock()
	defer task.Lock.Unlock()
	//dont do any operation if task is untouched
	if task.Status == model.FAILED {
		log.Printf("pushBack task with id %s with status %s", task.Id, task.Status)
		return true, false
	}
	//if task is not completed, then check if it has ended it time limit
	if  ((task.Status == model.TIMEOUT) || (task.Status == model.SUCCESS)) {
		log.Printf("remove task with id %s with status %s", task.Id, task.Status)
		return false, true
	}
	//log.Printf("retain task %s as it is still valid", task.Id)
	return false, false
}