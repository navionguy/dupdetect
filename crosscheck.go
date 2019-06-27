package main

import (
	"fmt"
)

// send all of my codes to one of my friends
// I don't need to wait on him to answer
// if there is a problem, he will tell mom
func sendToFriends(w *worker, fname string, c codeMap) {
	var rq request

	rq.fname = fname
	rq.codes = c

	w.inp <- rq
}

func listenForFriends(w *worker, keys codeMap) {
	fmt.Printf("%s listening for %d friends\n", w.fname, w.listenCount)

	for {
		if w.listenCount == 0 {
			break // I've heard from all my friends, so I'm done.
		}

		chk := <-w.inp

		fmt.Printf("asked to check %d codes from file %s\n", len(chk.codes), chk.fname)
		w.listenCount--

		go checkFriend(w, chk, keys)
	}
}

func checkFriend(w *worker, chk request, keys codeMap) {
	var res results

	for k := range chk.codes {
		_, ok := keys[k]

		if ok {
			res.codesSeen = -1
			res.resultsDescription = fmt.Sprintf("key value %s found in both %s and %s\n", k, w.fname, chk.fname)
			w.outp <- res
			return
		}
	}
	res.codesSeen = len(chk.codes)
	res.resultsDescription = checkOp
	w.outp <- res
	return
}
