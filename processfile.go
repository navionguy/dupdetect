package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const loadOp string = "Loaded"
const checkOp string = "CrossCheck"

// Worker holds control information for each worker process
type worker struct {
	fname       string
	outp        chan results
	inp         chan request
	friends     []*worker
	listenCount int
}

// codeMap is a map of all the codes a worker has loaded
type codeMap map[string]bool

// request is passed from one worker to the next checking for dupes
type request struct {
	fname string
	codes codeMap
}

// Results passes out what was learned in processing
type results struct {
	resultsDescription string
	codesSeen          int
}

// processFile takes a single file and loads it while check for duplicate entries
// He waits until all files have loaded and then starts his values against all
// other file values
func processFile(w *worker) {
	fmt.Printf("Processing file %s\n", w.fname)
	var res results

	f, err := os.Open(w.fname)

	if err != nil {
		res.codesSeen = -1
		res.resultsDescription = fmt.Sprintf("Unable to open file %s", w.fname)
		w.outp <- res
		return
	}

	keys, err := loadFile(f, w.fname)

	if err != nil {
		res.codesSeen = -1
		res.resultsDescription = fmt.Sprintf("Unable to process file %s, error: %s", w.fname, err.Error())
		w.outp <- res
		return
	}

	// we know have two things to do,
	// 1. Take codes from workers after us in the list and check if they collide with one of mine
	// 2. Check my codes against those held by workers before me in the list
	//
	// Two things to do, more goroutines!
	// start one to listen for other workers to send their code list

	go listenForFriends(w, keys)

	// now send my codes to all my friends to check against theirs
	// I don't actually need to know what they find

	for _, f := range w.friends {
		go sendToFriends(f, w.fname, keys)
	}

	// I've sent all my codes to my friends, time to close out my load operation

	res.resultsDescription = loadOp
	res.codesSeen = len(keys)
	w.outp <- res
}

// loads a file of codes and puts them into a map for easy searching later
// while loading, he checks for duplicates.  If one is found, the main
// process is notified and execution ends
func loadFile(r io.ReadCloser, fname string) (codeMap, error) {
	kf := csv.NewReader(r)
	defer r.Close()
	kf.FieldsPerRecord = -1 // however many are on a line is cool

	keys := make(codeMap)

	for {
		values, err := kf.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			return keys, err
		}

		for _, v := range values {
			v = strings.TrimLeft(v, " ") // whitespace don't matter
			v = strings.TrimRight(v, " ")
			_, ok := keys[v]

			if ok {
				return keys, errors.New("duplicate found")
			}
			//fmt.Printf("%s added %s\n", fname, v)
			keys[v] = true
		}
	}
	return keys, nil
}
