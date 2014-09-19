package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"os"
	"path"

	"github.com/jvlmdr/go-cv/detect"
	"github.com/jvlmdr/go-cv/feat"
	_ "github.com/jvlmdr/go-cv/hog"
	"github.com/jvlmdr/go-cv/imsamp"
	"github.com/jvlmdr/go-file/fileutil"
	"github.com/nfnt/resize"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s [flags] tmpl.(gob|json) transform.json\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var (
		// Dataset options.
		name = flag.String("dataset", "usatest", "Dataset identifier.")
		dir  = flag.String("dir", "", "Location of dataset. Empty means working dir.")
		// Detection options.
		pyrStep  = flag.Float64("pyr-step", 1.2, "Geometric scale steps in image pyramid.")
		maxScale = flag.Float64("max-scale", 1, "Maximum amount to scale image. Greater than 1 is upsampling.")
		maxIOU   = flag.Float64("max-iou", 0.3, "Maximum IOU that two detections can have before NMS.")
		margin   = flag.Int("margin", 0, "Spatial bin parameter to HOG.")
		// Validation options.
		minValIOU      = flag.Float64("min-val-iou", 0.5, "Minimum IOU for a detection to be validated.")
		minIgnoreCover = flag.Float64("min-ignore-cover", 0, "Minimum that a detection must be covered by an ignore region to be ignored.")
		// Display options.
		numShow = flag.Int("num-show", 4, "Number of detections to show per image")
	)
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	var (
		tmplFile      = flag.Arg(0)
		transformFile = flag.Arg(1)
	)

	// Get dataset from name.
	dataset, err := datasetByName(*name)
	if err != nil {
		log.Fatalln(err)
	}
	// Load template from file.
	var tmpl *detect.FeatTmpl
	if err := fileutil.LoadExt(tmplFile, &tmpl); err != nil {
		log.Fatalln("load template:", err)
	}
	// Load transform from file.
	var transform *feat.Marshaler
	if err := fileutil.LoadJSON(transformFile, &transform); err != nil {
		log.Fatalln("load transform:", err)
	}

	overlap := func(a, b image.Rectangle) bool { return detect.IOU(a, b) > *maxIOU }
	opts := detect.MultiScaleOpts{
		MaxScale:    *maxScale,
		PyrStep:     *pyrStep,
		Interp:      resize.Bicubic,
		Transform:   transform,
		Pad:         feat.Pad{feat.UniformMargin(*margin), imsamp.Continue},
		DetFilter:   detect.DetFilter{LocalMax: true, MinScore: math.Inf(-1)},
		SupprFilter: detect.SupprFilter{MaxNum: 0, Overlap: overlap},
	}

	var val *detect.ValSet
	err = fileutil.Cache(&val, "val-set.json", func() (*detect.ValSet, error) {
		return testAll(dataset, *dir, tmpl, opts, *minValIOU, *minIgnoreCover, *numShow)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func testAll(dataset *Dataset, dir string, tmpl *detect.FeatTmpl, opts detect.MultiScaleOpts, minValIOU, minIgnoreCover float64, numShow int) (*detect.ValSet, error) {
	// Load each image and perform multi-scale detection.
	rootDir := path.Join(dir, "data-"+dataset.Dir)
	var vals []*detect.ValSet
	for set, seqs := range dataset.Seqs {
		for _, seq := range seqs {
			relDir := path.Join(fmt.Sprintf("set%02d", set), fmt.Sprintf("V%03d", seq))
			// Check that image directory exists.
			imDir := path.Join(rootDir, "images", relDir)
			if err := isDir(imDir); err != nil {
				return nil, fmt.Errorf("check image dir: %v", err)
			}
			// Check that annotation directory exists.
			annotDir := path.Join(rootDir, "annotations", relDir)
			if err := isDir(annotDir); err != nil {
				return nil, fmt.Errorf("check annotation dir: %v", err)
			}

			// Create results directory.
			resDir := path.Join("res", relDir)
			if err := os.MkdirAll(resDir, 0755); err != nil {
				return nil, fmt.Errorf("create results dir: %v", err)
			}
			// Create visualizations directory.
			visDir := path.Join("vis-dets", relDir)
			if err := os.MkdirAll(visDir, 0755); err != nil {
				return nil, fmt.Errorf("create visualizations dir: %v", err)
			}

			var visFiles []string
			for j := 0; ; j++ {
				frame := (j+1)*dataset.Skip - 1
				imFile := path.Join(imDir, fmt.Sprintf("I%05d.%s", frame, dataset.Ext))
				// Continue until image file does not exist.
				if _, err := os.Stat(imFile); os.IsNotExist(err) {
					break
				}
				log.Printf("detect: seq %s, frame %d", relDir, frame)
				// Load image file.
				im, err := loadImage(imFile)
				if err != nil {
					return nil, err
				}
				// Perform multi-scale detection.
				dets := detect.MultiScale(im, tmpl, opts)
				// Save detections for each image to file.
				resFile := path.Join(resDir, fmt.Sprintf("I%05d.txt", frame))
				if err := saveDets(resFile, dets); err != nil {
					return nil, fmt.Errorf("save detections: %v", err)
				}
				// Load annotations and validate detections.
				annotFile := path.Join(annotDir, fmt.Sprintf("I%05d.txt", frame))
				annot, err := loadAnnot(annotFile)
				if err != nil {
					return nil, fmt.Errorf("load annotations: %v", err)
				}
				refs, ignore := Rects(annot, Reasonable())
				val := detect.Validate(dets, refs, ignore, minValIOU, minIgnoreCover)
				vals = append(vals, val.Set())
				// Save visualization of detections.
				visFile := path.Join(visDir, fmt.Sprintf("I%05d.jpg", frame))
				if err := visDets(visFile, imFile, val, numShow); err != nil {
					return nil, fmt.Errorf("visualize detections: %v", err)
				}
				visFiles = append(visFiles, fmt.Sprintf("I%05d.jpg", frame))
			}

			// Create index of visualizations.
			visIndexFile := path.Join(visDir, "index.html")
			if err := saveVisIndex(visIndexFile, visFiles); err != nil {
				return nil, fmt.Errorf("save index of visualizations: $v", err)
			}
		}
	}
	return detect.MergeValSets(vals...), nil
}

func isDir(fname string) error {
	if info, err := os.Stat(fname); err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("is not dir: %s", fname)
	}
	return nil
}

func loadImage(fname string) (image.Image, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("load image: %s: %v", fname, err)
	}
	return im, nil
}

func saveDets(fname string, dets []detect.Det) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()
	return writeDets(file, dets)
}

func writeDets(w io.Writer, dets []detect.Det) error {
	for _, det := range dets {
		r := det.Rect
		fmt.Fprintf(w, "%d,%d,%d,%d,%e\n", r.Min.X, r.Min.Y, r.Dx(), r.Dy(), det.Score)
	}
	return nil
}
