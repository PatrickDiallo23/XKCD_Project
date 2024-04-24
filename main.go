package main

import (
	"flag"
	"fmt"
	"xkcdcomics/cmd/task1"
	"xkcdcomics/cmd/task2"
)

func main() {
	fmt.Println("Hello, World! Welcome to XKCD go bootcamp")

	mode := flag.String("mode", "", "mode of retrieval (sequential, concurrent or xckd)")
	flag.Parse()

	if *mode == "sequential" || *mode == "concurrent" {
		task1.RunTask1(mode)
	} else if *mode == "xkcd" {
		task2.RunTask2()
	} else {
		fmt.Println("Choose a valid flag, please...")
		fmt.Println("sequential, concurrent or xkcd")
	}
}
