package main

import (
	"fmt"
)

type ThreadPool struct {
	numThreads int
	inputs chan threadPoolTask
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


type threadPoolTask interface {
	say(x int) string
	equal(x, y int) bool
}



func main() {
	tp := ThreadPool{}
	tp.init(1)

	tp.start()
	tp.inputs <- func() {
		say(1)
	}
	tp.inputs <- func() {
		otherTask(2)
	}

}