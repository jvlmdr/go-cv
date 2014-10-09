package convfeat_test

import (
	"image"
	"testing"

	"github.com/jvlmdr/go-cv/convfeat"
	"github.com/jvlmdr/go-cv/feat"
)

func TestMarshaler(t *testing.T) {
	transforms := []feat.RealMarshalable{
		new(convfeat.PosPart),
		new(convfeat.PosNegPart),
		new(convfeat.PosNegPart),
		new(convfeat.IsPos),
		new(convfeat.Sign),
		&convfeat.MaxPool{image.Pt(3, 4), 2},
		&convfeat.SumPool{image.Pt(3, 4), 2},
		&convfeat.AdjChanNorm{5, 2, 1e-4, 0.75},
	}
	for _, phi := range transforms {
		err := feat.TestRealMarshaler(phi)
		if err != nil {
			t.Error(err)
		}
	}
}
