package inria

import (
	"bufio"
	"io"
	"os"
)

func loadLines(fname string) ([]string, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return readLinesFrom(file)
}

func readLinesFrom(r io.Reader) ([]string, error) {
	var lines []string
	// Read through lines of file.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}
