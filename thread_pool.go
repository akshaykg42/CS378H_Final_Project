package main

import (
	"fmt"
	"strconv"
)

type ThreadPool struct {
	numThreads int
	inputs chan func() interface{}
	outputs chan interface{}
}

func (tp *ThreadPool) init(size int, inputsChannelBufferSize int, outputsChannelBufferSize int) {
	tp.numThreads = size
	tp.inputs = make(chan func() interface{}, inputsChannelBufferSize)
	tp.outputs = make(chan interface{}, outputsChannelBufferSize)
}

func (tp *ThreadPool) runThread(id int) {
	for {
		fn, more := <- tp.inputs
		if !more {
			break
		}
		out := fn()
		tp.outputs <- out
	}
}

func (tp *ThreadPool) start() {
	for i := 0; i < tp.numThreads; i++ {
		go tp.runThread(i)
	}
}

func (tp *ThreadPool) finish() {
	close(tp.inputs)
	close(tp.outputs)
}


/////////////////////

func say(x int) string {
	fmt.Printf("Input is: %d\n", x)
	return strconv.Itoa(x)
}

func add(x int, y int) int {
	fmt.Printf("%d + %d = %d\n", x, y, x + y)
	return x + y
}

func main() {
	// Creates and initializes a thread pool of size 2, with input and
	// output channels with buffer size 0. Then, starts the threads.
	tp := ThreadPool{}
	tp.init(2, 0, 0)
	tp.start()

	// This chunk of code uses a loop over a select to asynchronously
	// delegate tasks to the thread pool. In this situation, the thread
	// pool runs say() and add() five times each.
	addRunCount := 0
	sayRunCount := 0
	addReceiveCount := 0
	sayReceiveCount := 0

	// this is possibly the most disgusting code i have ever written
	for addRunCount < 6 || sayRunCount < 6 || addReceiveCount < 6 || sayReceiveCount < 6{
		if addRunCount < 6 {
			i := addRunCount
			select {
			case tp.inputs <- func() interface{} { return add(i, i) }:
				addRunCount++
			default:
			}
		}

		if sayRunCount < 6 {
			i := sayRunCount
			select {
			case tp.inputs <- func() interface{} { return say(i) }:
				sayRunCount++
			default:
			}
		}

		if addReceiveCount < 6 || sayReceiveCount < 6{
			select {
			case output := <- tp.outputs:
				switch output.(type) {
				case string:
					sayReceiveCount++
					fmt.Printf("Outputted a String: %s\n", output)
				case int:
					addReceiveCount++
					fmt.Printf("Outputted an Int: %d\n", output)
				}
			default:
			}
		}
	}

	// Close the input and output channels, killing the threads.
	tp.finish()
}