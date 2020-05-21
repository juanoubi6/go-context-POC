package main

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wgGorutines sync.WaitGroup

//Expected result: the gorutine cancellation order should be
// cancelCtx and cancelCtxSonCtx at the same time
// timeoutCtx
// deadlineCtx
func main() {
	println("Number of processors: " + strconv.Itoa(runtime.NumCPU()))
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Create new empty context
	ctx := context.Background()

	// Create a WaitGroup to wait for all gorutines to finish before exiting
	wgGorutines.Add(4)

	// Create 4 new contexts:
	// Cancel context: can be cancelled from outside with the cancel() function
	// Timeout context: after a certain amount of time, the context is cancelled automatically. We are going to create
	//                  an additional timeOut context that inherits from a cancel context.
	// Deadline context: when the clock reaches a certain time, the context is cancelled automatically. Useful
	//                   to propagate timeout across APIs.
	cancelCtx, cancel := context.WithCancel(ctx)
	timeoutCtx, _ := context.WithTimeout(ctx, time.Second*5)
	cancelCtxSonCtx, _ := context.WithTimeout(cancelCtx, time.Second*10)
	deadlineCtx, _ := context.WithDeadline(ctx, time.Now().Add(time.Second*7))

	go worker(cancelCtx, "Closing cancelCtx gorutine")
	go worker(cancelCtxSonCtx, "Closing cancelCtxSonCtx gorutine")
	go worker(timeoutCtx, "Closing timeoutCtx gorutine")
	go worker(deadlineCtx, "Closing deadlineCtx gorutine")

	// Cancel the cancelCtx goroutine and any context that inherits from it
	time.Sleep(time.Second * 2)
	cancel()

	wgGorutines.Wait()
	println("Number of gorutines executing: " + strconv.Itoa(runtime.NumGoroutine()))

}

func worker(ctx context.Context, byeMsg string) {
	defer wgGorutines.Done()

	for {
		select {
		case <-ctx.Done():
			println(byeMsg)
			return
		default:
			println("Doing some work in gorutine")
			time.Sleep(time.Second)
		}
	}
}
