package main

import (
	"fmt"
	"os"
	"strings"
)

// scan a set of files to make sure none of the codes are duplicated, either within the file
// or across the set of files
func main() {
	fmt.Println("dupdetect")

	files := os.Args[1:]      // assume that every parameter is a filename
	res := make(chan results) // either a count of codes processed or failure code

	var workers []*worker

	// build an array of workers, one per file
	// each worker needs to know about all the workers created ahead of ti

	for i, file := range files {
		w := new(worker)
		w.fname = file
		w.outp = res
		w.inp = make(chan request)
		w.friends = workers[:i]

		workers = append(workers, w)
	}

	// launch each of the workers

	totalThreads := len(workers)
	for i, w := range workers {
		w.listenCount = len(workers) - (i + 1)
		totalThreads += w.listenCount
		go processFile(w)
	}

	// now wait for the result and print the results
	// each file processed will result in two result packets
	// the first is sent when the file has fully loaded and reports the number of codes loaded without a duplicate
	// the second is number of codes checked for his "friend" workers
	// if at anytime, the resutt comes back -1, a duplicate has been found and we should just stop

	totalCodesSeen := 0
	totalCodesCross := 0

	for i := 0; i < totalThreads; i++ {
		result := <-res

		if result.codesSeen == -1 {
			// time to go, print result and get out
			fmt.Println(result.resultsDescription)
			return
		}

		if strings.Compare(result.resultsDescription, loadOp) == 0 {
			totalCodesSeen += result.codesSeen
			fmt.Printf("Saw %d loaded\n", result.codesSeen)
		}

		if strings.Compare(result.resultsDescription, checkOp) == 0 {
			totalCodesCross += result.codesSeen
			fmt.Printf("Saw %d codes checked\n", result.codesSeen)
		}
	}

	fmt.Printf("A total of %d codes checked, %d codes cross checked, no duplicates found.\n", totalCodesSeen, totalCodesCross)
}
