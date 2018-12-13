// [_Command-line flags_](http://en.wikipedia.org/wiki/Command-line_interface#Command-line_option)
// are a common way to specify options for command-line
// programs. For example, in `wc -l` the `-l` is a
// command-line flag.

package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//XORs 2 strings
func xor(s1, s2 string) (output string) {
	output = ""
	 for i := range s1 {
	 	x,_ := strconv.Atoi(s1[i : i + 1])
	 	y,_ := strconv.Atoi(s2[i : i + 1])
	 	output += strconv.Itoa(x ^ y)
	 }
	 //fmt.Println(output)
	 return output
}

func prf(data int, blocksize int) (out string) {
	out = ""
    for i := 0; i < blocksize; i++ {
        out = out + strconv.Itoa(rand.Intn(2))
    }
	//fmt.Printf("PRF output was: %s\n", out)
	return out
}

//Takes input string and block size, returns the string separated into blocks
//Note: Assumes that len(string) % blocksize == 0
func get_blocks(data string, blocksize int) (blocks []string){
    for i := 0; i < len(data) / blocksize; i++ {
    	//block, _ := strconv.Atoi(data[blocksize * i : blocksize * (i + 1)])
        blocks = append(blocks, data[blocksize * i : blocksize * (i + 1)])
    }

    return blocks
}

//Takes plaintext and returns encrypted text
func crypto_ctr(plaintext string, blocksize int, iv int) (output string){
    blocks := get_blocks(plaintext, blocksize)

    for i, block := range(blocks){
        output += xor(block, prf(iv + i, blocksize))
    }
    return output
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	var input string
    flag.StringVar(&input, "input", "ctr_input.txt", "String-valued path to an input file")
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

    start := time.Now()

    crypto_ctr(plaintext, blocksize, iv)

    t := time.Now()
    elapsed := t.Sub(start)
    //Microseconds
    fmt.Println(float64(elapsed.Nanoseconds()) / 1000000.0)

    //fmt.Println(output)
}
