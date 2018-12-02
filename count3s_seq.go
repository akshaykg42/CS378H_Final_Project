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

func count3s(array []int64){
	count := 0

	for _, num := range(array){
		if(num == 3){
			count++
		}
	}

	//fmt.Println(count)
}

func main() {
    var input string
    flag.StringVar(&input, "input", "count3s_input.txt", "String-valued path to an input file")
    flag.Parse()

    var data []int64

    file, _ := os.Open(input)

    defer file.Close()

    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
    	num, _ := strconv.ParseInt(scanner.Text(), 10, 64)
    	data = append(data, num)
    }

    start := time.Now()

    count3s(data)

    t := time.Now()
    elapsed := t.Sub(start)
    fmt.Println(float64(elapsed.Nanoseconds()) / 1000000.0)
}