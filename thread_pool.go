package main

import (
	"fmt"
)

type ThreadPool struct {
	numThreads int
	inputs chan func()
	outputs chan int
}

func (tp *ThreadPool) init(size int) {
	tp.numThreads = size
	tp.inputs = make(chan func())
	tp.outputs = make(chan int)
}

func (tp *ThreadPool) runThread(id int) {
	for {
		fn, more := <- tp.inputs
		if !more {
			break
		}
		fn()
	}
}

func (tp *ThreadPool) start() {
	for i := 0; i < tp.numThreads; i++ {
		go tp.runThread(i)
	}
}


/////////////////////


func task(x int) {
	fmt.Printf("x is: %d\n", x)
}

func otherTask(y int)  {
	fmt.Printf("y is: %d\n", y)
}

func main() {
	tp := ThreadPool{}
	tp.init(1)

	tp.start()
	tp.inputs <- func() {
		task(1)
	}
	tp.inputs <- func() {
		otherTask(2)
	}

}