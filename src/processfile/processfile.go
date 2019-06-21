package processfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ProcessFile takes a single file and loads it while check for duplicate entries
// He waits until all files have loaded and then starts his values against all
// other file values
func ProcessFile(fname string, outp chan int) {
	f, err := os.Open(fname)

	if err != nil {
		return
	}

	ct, err := loadFile(f)

	if err != nil {
		outp <- -1
		return
	}

	outp <- ct
}

func loadFile(r io.ReadCloser) (int, error) {
	kf := csv.NewReader(r)
	defer r.Close()
	kf.FieldsPerRecord = -1 // however many are on a line is cool

	keys := make(map[string]bool)

	for {
		values, err := kf.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("error %s\n", err.Error())
			break
		}

		for _, v := range values {
			v = strings.TrimLeft(v, " ") // whitespace don't matter
			v = strings.TrimRight(v, " ")
			_, ok := keys[v]

			if ok {
				return len(keys), errors.New("duplicate found")
			}
			keys[v] = true
		}
	}
	return len(keys), nil
}
