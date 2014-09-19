package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"

	"github.com/jvlmdr/go-cv/detect"
)

func visDets(dst, src string, val *detect.ValImage, n int) error {
	var args []string
	args = append(args, src)
	args = append(args, "-strokewidth", "2")
	args = append(args, "-undercolor", "#00000080")
	// Annotate highest scoring annotations.
	for i, det := range val.Dets {
		if i >= n && !det.True {
			continue
		}
		// Draw rectangles.
		args = append(args, "-fill", "none")
		if det.True {
			args = append(args, "-stroke", "green")
		} else {
			args = append(args, "-stroke", "yellow")
		}
		args = append(args, "-draw", rectStr(det.Rect))
		if det.True {
			// Draw matched reference.
			args = append(args, "-stroke", "blue")
			args = append(args, "-draw", rectStr(det.Ref))
		}
		// Label with score.
		args = append(args, "-fill", "white", "-stroke", "none")
		pos := fmt.Sprintf("+%d+%d", det.Rect.Min.X, det.Rect.Max.Y)
		text := fmt.Sprintf("%.4g (%d)", det.Score, i+1)
		args = append(args, "-annotate", pos, text)
	}
	// Annotate detections which were completely missed.
	args = append(args, "-fill", "none", "-stroke", "red")
	for _, rect := range val.Misses {
		args = append(args, "-draw", rectStr(rect))
	}
	args = append(args, dst)
	cmd := exec.Command("convert", args...)
	return cmd.Run()
}

func rectStr(r image.Rectangle) string {
	return fmt.Sprintf("rectangle %d,%d %d,%d", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

func saveVisIndex(fname string, ims []string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeVisIndex(file, ims)
}

func writeVisIndex(w io.Writer, ims []string) error {
	for _, im := range ims {
		fmt.Fprintf(w, "<div><img src=\"%s\" /></div>\n", im)
	}
	return nil
}
