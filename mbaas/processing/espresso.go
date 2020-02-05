package processing

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type factor int

const (
	dontcare factor = 0
	pos      factor = 1
	neg      factor = 2
)

type expression [][]factor

func espresso(seq sequence) (expression, error) {
	// Generate Espresso file.
	infile, err := ioutil.TempFile("", "*.espresso")
	if err != nil {
		return nil, err
	}
	//defer os.Remove(infile.Name())
	if err = writeInputFile(infile, seq); err != nil {
		return nil, err
	}

	// Run Espresso and read output.
	espresso, err := exec.Command("espresso", infile.Name()).Output()
	if err != nil {
		return nil, err
	}

	// Process Espresso output.
	return convertOutput(strings.Split(string(espresso), "\n"))
}

func writeInputFile(f *os.File, seq sequence) error {
	nbits := numberOfBits(len(seq) - 1)

	fmt.Fprintf(f, ".i %d\n", nbits)
	fmt.Fprintf(f, ".o 1\n")

	for i, b := range seq {
		if _, err := fmt.Fprintf(f, "%0*b %s\n", nbits, i, b.String()); err != nil {
			return err
		}
	}

	return f.Close()
}

func convertOutput(lines []string) (expression, error) {
	result := expression{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0:1] == "." {
			continue
		}
		bits := strings.Fields(line)[0]
		term := []factor{}
		for _, ch := range bits {
			switch ch {
			case '-':
				term = append(term, dontcare)
			case '0':
				term = append(term, neg)
			case '1':
				term = append(term, pos)
			}
		}
		result = append(result, term)
	}
	return result, nil
}
