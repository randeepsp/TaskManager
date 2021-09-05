package adders

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"model"
	"os"
	"time"
	"queue"
)

//starts adding the tasks
func StartTaskAdder(ctx context.Context, queue *queue.Queue, endProgram chan os.Signal){
	log.Print("starting to add tasks")
	ended := false
	for !ended {
		select {
		//end the program when a signal is recieved on channel
		case <-endProgram :
				log.Print("Stopped adding tasks into queue")
				ended = true
				break
		default:
			//a dummy random task object,can be replaced with actual content
			taskData := "randomTask"
			uuid := addTask(ctx, queue, taskData)
			log.Printf("Added task %s to queue", uuid)
		}
	}
	log.Printf("Completed adding tasks")
}

//created a new task based on taskData and adds to the queue
func addTask(ctx context.Context, queue *queue.Queue, taskData string) string{
	uuid := uuidGenerator()
	task := model.Task{
		Id:           uuid,
		Status:       model.UNTOUCHED,
		CreationTime: time.Now(),
		TaskData:     taskData,
	}
	queue.Enqueue(&task)
	return uuid
}

//returns a random uuid string
func uuidGenerator() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(u)
}