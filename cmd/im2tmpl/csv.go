package main

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/jvlmdr/go-cv/rimg64"
)

func loadImageCSV(fname string) (*rimg64.Multi, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return readImageCSV(file)
}

func readImageCSV(r io.Reader) (*rimg64.Multi, error) {
	cr := csv.NewReader(r)
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	return parseImageCSV(rows)
}

func parseImageCSV(rows [][]string) (*rimg64.Multi, error) {
	// Convert to numbers.
	var (
		u, v, w []int
		f       []float64
	)
	for _, r := range rows {
		if len(r) == 0 {
			continue
		}
		if len(r) != 4 {
			panic("wrong number of elements in row")
		}
		ui, err := strconv.ParseInt(r[0], 10, 32)
		if err != nil {
			return nil, err
		}
		vi, err := strconv.ParseInt(r[1], 10, 32)
		if err != nil {
			return nil, err
		}
		wi, err := strconv.ParseInt(r[2], 10, 32)
		if err != nil {
			return nil, err
		}
		fi, err := strconv.ParseFloat(r[3], 64)
		if err != nil {
			return nil, err
		}
		u = append(u, int(ui))
		v = append(v, int(vi))
		w = append(w, int(wi))
		f = append(f, fi)
	}

	// Take max over u, v, w.
	var width, height, channels int
	for i := range u {
		width = max(u[i]+1, width)
		height = max(v[i]+1, height)
		channels = max(w[i]+1, channels)
	}

	// Set pixels in image.
	im := rimg64.NewMulti(width, height, channels)
	for i := range u {
		im.Set(u[i], v[i], w[i], f[i])
	}
	return im, nil
}

func saveImageCSV(fname string, im *rimg64.Multi) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeImageCSV(file, im)
}

func writeImageCSV(w io.Writer, im *rimg64.Multi) error {
	cw := csv.NewWriter(w)
	rows := formatImageCSV(im)
	return cw.WriteAll(rows)
}

func formatImageCSV(im *rimg64.Multi) [][]string {
	var rows [][]string
	for u := 0; u < im.Width; u++ {
		for v := 0; v < im.Height; v++ {
			for w := 0; w < im.Channels; w++ {
				r := make([]string, 4)
				r[0] = strconv.FormatInt(int64(u), 10)
				r[1] = strconv.FormatInt(int64(v), 10)
				r[2] = strconv.FormatInt(int64(w), 10)
				r[3] = strconv.FormatFloat(im.At(u, v, w), 'g', -1, 64)
				rows = append(rows, r)
			}
		}
	}
	return rows
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
