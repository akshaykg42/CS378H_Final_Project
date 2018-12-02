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
    "time"
)

//XORs 2 strings
func xor(s1, s2 string) (output string) {
        for i := 0; i < len(s1); i++ {
                output += string(s1[i] ^ s2[i % len(s2)])
        }

        return output
}

//Currently just outputs 'x' * blocksize, needs to be an actual encryption function
func encrypt(data int, key int, blocksize int) (out string) {
    for i := 0; i < blocksize; i++ {
        out += "x"
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

//Takes plaintext and returns encrypted text (or vice versa)
func crypto_ctr(plaintext string, key int, blocksize int, iv int) (output string){
    blocks := get_blocks(plaintext, blocksize)

    for i, block := range(blocks){
        output += xor(block, encrypt(iv + i, key, blocksize))
    }

    return output
}

func main() {
    var input string
    flag.StringVar(&input, "input", "ctr_input.txt", "String-valued path to an input file")
    flag.Parse()

    var plaintext string
    key := 1
    blocksize := 1
    iv := 1

    file, _ := os.Open(input)

    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
    	plaintext = scanner.Text()
    }

    start := time.Now()

    crypto_ctr(plaintext, key, blocksize, iv)

    t := time.Now()
    elapsed := t.Sub(start)
    fmt.Println(float64(elapsed.Nanoseconds()) / 1000000.0)
}