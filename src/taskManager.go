package main

import (
	"adders"
	"cleaners"
	"context"
	"executors"
	"log"
	"model"
	"os"
	"os/signal"
	"queue"
	"sync"
	"time"
)

func main() {
	log.Print("Started task management")
	var              wg sync.WaitGroup
	model.Polltime   = time.Millisecond*15
	ctx              := context.Background()
	taskQ            := queue.NewQueue()
	//handle Ctrl+C signal to end adding tasks
	endProgramSignal := make(chan os.Signal, 1)
	signal.Notify(endProgramSignal, os.Interrupt)

	//program will wait for the cleaner to complete as it can be ended once all tasks are removed
	//not really using it right now
	wg.Add(1)
	//start adding tasks
	go adders.StartTaskAdder(ctx, taskQ, endProgramSignal)
	//start executor
	go executors.StartExecutor(ctx, taskQ)
	//start cleaner
	go cleaners.StartCleaner(ctx, &wg, taskQ)
	wg.Wait()
	log.Print("Task management completed")
}