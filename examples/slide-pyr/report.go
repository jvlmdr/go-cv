package main

import (
	"github.com/jackvalmadre/go-cv"
	"github.com/jackvalmadre/go-cv/imgpyr"

	"bufio"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"os"
)

var reportTmpl string = `<html>
<head>
<style type="text/css">
table {
	border-collapse: collapse;
}
table, th, td {
	border: 1px solid black;
}
th, td {
	padding: 1ex;
}
th.input, td.input {
	background-color: #cccccc;
}
</style>
</head>
<body>

<table>
	<tr>
		<th>Rank</th>
		<th>Image</th>
		<th>Score</th>
		<th>Level</th>
	</tr>
	{{range .}}
	<tr>
		<td>Rank</td>
		<td><img src="{{.Image}}" /></td>
		<td>{{.Score}}</td>
		<td>{{.Level}}</td>
	</tr>
	{{end}}
</table>

</body>
</html>
`

type Row struct {
	Rank  int
	Image string
	Score float64
	Level int
}

func report(reportFile string, pyr *imgpyr.Pyramid, dets []PyrPos, size image.Point, scores []cv.RealImage, featRate int) {
	rows := make([]Row, len(dets))
	for i, det := range dets {
		img := extractImage(pyr, det, size, featRate)
		// Save to file.
		rows[i].Image = fmt.Sprintf("%06d.jpg", i)
		if err := saveImage(rows[i].Image, img); err != nil {
			log.Fatalln("could not save detection image:", err)
		}
		// Other details.
		rows[i].Rank = i + 1
		rows[i].Score = scores[det.Level].At(det.Pos.X, det.Pos.Y)
		rows[i].Level = det.Level
	}

	// Generate HTML.
	tmpl, err := template.New("detections").Parse(reportTmpl)
	if err != nil {
		log.Fatalln("could not parse template:", err)
	}
	if err := executeSave(reportFile, tmpl, rows); err != nil {
		log.Fatalln("could not execute template to file:", err)
	}
}

func executeSave(filename string, tmpl *template.Template, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	buf := bufio.NewWriter(file)
	defer buf.Flush()
	return tmpl.Execute(buf, data)
}

type subImager interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func extractImage(pyr *imgpyr.Pyramid, det PyrPos, size image.Point, featRate int) image.Image {
	r := image.Rectangle{det.Pos, det.Pos.Add(size)}
	// Scale rectangle from feature image back to real image.
	r = scaleRect(float64(featRate), r)
	// Try to upgrade pyramid level to a subImager.
	img, ok := pyr.Levels[det.Level].(subImager)
	if !ok {
		panic("pyramid level is not a subImager")
	}
	// Extract sub-image.
	return img.SubImage(r)
}

func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		return err
	}
	return nil
}
