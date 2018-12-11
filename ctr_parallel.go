package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
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

func xor(s1, s2 string) string {
	output := ""
	for i := range s1 {
		x,_ := strconv.Atoi(s1[i : i + 1])
		y,_ := strconv.Atoi(s2[i : i + 1])
		output += strconv.Itoa(x ^ y)
	}
	return output
}

//Currently just outputs 'x' * blocksize, needs to be an actual encryption function
func prf(data int, blocksize int) (out string) {
	out = ""
	for i := 0; i < blocksize; i++ {
		out = out + strconv.Itoa(rand.Intn(2))
	}
	return out
}

//Takes input string and block size, returns the string separated into blocks
//Note: Assumes that len(string) % blocksize == 0
func get_blocks(data string, blocksize int) (blocks []string){
	for i := 0; i < len(data) / blocksize; i++ {
		blocks = append(blocks, data[blocksize * i : blocksize * (i + 1)])
	}
	return blocks
}

func ctr_parallelized(plaintext string, blocksize int, iv int) string{
	blocks := get_blocks(plaintext, blocksize)

	tp := ThreadPool{}
	tp.init(1, 0, 0)
	tp.start()

	inputCount := 0
	outputCount := 0
	output := make([]string, len(blocks))
	for inputCount < len(blocks) || outputCount < len(blocks){
		if inputCount < len(blocks) {
			i := inputCount
			s := blocks[inputCount]
			select {
			case tp.inputs <- func() interface{} {
				out := [2]string{xor(s, prf(iv + i, blocksize)), strconv.Itoa(i)}
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
				case [2]string:
					outputCount++
					i, _ := strconv.Atoi(result[1])
					output[i] = result[0]
				}
			default:
			}
		}
	}

	tp.finish()

	return strings.Join(output, "")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var input string
	var workers int
	flag.StringVar(&input, "input", "ctr_input.txt", "String-valued path to an input file")
	flag.IntVar(&workers, "workers", 1, "Number of workers to spawn")
	flag.Parse()
	var plaintext string
	blocksize := 0
	iv := 0

	file, _ := os.Open(input)
	defer file.Close()
	reader := bufio.NewReader(file)

	line, _ := reader.ReadBytes('\n')
	line = line[:len(line) - 1]
	blocksize, _ = strconv.Atoi(string(line))
	line, _ = reader.ReadBytes('\n')
	line = line[:len(line) - 1]
	plaintext = string(line)

	fmt.Println(ctr_parallelized(plaintext, blocksize, iv))

}
