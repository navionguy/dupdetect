package main

import (
	"dupdetect/src/processfile"
	"fmt"
	"os"
)

type worker struct {
	fname   string
	outp    chan int
	friends []chan string
}

func main() {
	fmt.Println("dupdetect")

	files := os.Args[1:]
	var workers []worker

	// launch the workers to load their files

	for _, file := range files {
		fmt.Printf("Processing file %s\n", file)

		w := worker{file, make(chan int), nil}
		workers = append(workers, w)
		go processfile.ProcessFile(file, w.outp)
	}

	// now wait till everybody finishes loading their files

	count := 0
	for _, w := range workers {
		ct := <-w.outp

		if ct == -1 {
			fmt.Printf("duplicate found in file %s\n", w.fname)
			return
		}
		count += ct
	}

	fmt.Printf("a total of %d codes were loaded\n", count)
}
