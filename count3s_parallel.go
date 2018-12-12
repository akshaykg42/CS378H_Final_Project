// [_Command-line flags_](http://en.wikipedia.org/wiki/Command-line_interface#Command-line_option)
// are a common way to specify options for command-line
// programs. For example, in `wc -l` the `-l` is a
// command-line flag.

package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
    "flag"
    "fmt"
    "os"
    "bufio"
    "strconv"
    "time"
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

//Splits array into a bunch of chunks of size blocksize (assuming it divides len(data))
func get_blocks(data []int64, blocksize int) (blocks [][]int64){
    for i := 0; i < len(data) / blocksize; i++ {
        blocks = append(blocks, data[blocksize * i : blocksize * (i + 1)])
    }
    return blocks
}

//Sequential count3s algorithm for each worker thread to run
func count3s(array []int64) int{
    count := 0

    for _, num := range(array){
        if(num == 3){
            count++
        }
    }

    return count
}

func count3s_parallel(array []int64, blocksize int64) int{
    //Get the blocks
    blocks := get_blocks(array, blocksize)

    //Initialize total count
    totalCount := 0

    tp := ThreadPool{}
    tp.init(1, 0, 0)
    tp.start()

    inputCount := 0
    outputCount := 0

    for inputCount < len(blocks) || outputCount < len(blocks){
        if inputCount < len(blocks) {
            block := blocks[inputCount]
            select {
            case tp.inputs <- func() interface{} {
                //Count 3s in this block and return in int array of size 1
                out := [1]int{count3s(block)}
                return out
            }:
                inputCount++
            default:
            }
        }

        if outputCount < len(blocks) {
            //fmt.Printf("output. len: %d count: %d\n", len(blocks), outputCount)
            select {
            case result := <- tp.outputs:
                switch result := result.(type) {
                //Wait for int array of size 1
                case [1]int:
                    outputCount++
                    //Add the number to the total count
                    totalCount += result[0]
                }
            default:
            }
        }
    }

    tp.finish()

    //Return total count
    return totalCount
}

func main() {
    //Path to Input
    var input string
    flag.StringVar(&input, "input", "count3s_input.txt", "String-valued path to an input file")
    flag.Parse()

    //Array containing all input numbers
    var data []int64

    file, _ := os.Open(input)

    defer file.Close()

    scanner := bufio.NewScanner(file)

    //First line of input is blocksize
    scanner.Scan()
    blocksize, _ := strconv.ParseInt(scanner.Text(), 10, 64)

    //Read in each number
    for scanner.Scan() {
    	num, _ := strconv.ParseInt(scanner.Text(), 10, 64)
    	data = append(data, num)
    }

    //Measuring Time
    start := time.Now()

    //Calling Parallel count3s
    count3s_parallel(data, blocksize)

    t := time.Now()
    elapsed := t.Sub(start)
    //Printing Time
    fmt.Println(float64(elapsed.Nanoseconds()) / 1000000.0)
}