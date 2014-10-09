package feat_test

import (
	"testing"

	"github.com/jvlmdr/go-cv/feat"
)

func TestImageMarshaler(t *testing.T) {
	var err error
	err = feat.TestImageMarshaler(new(feat.RGB))
	if err != nil {
		t.Error(err)
	}
	err = feat.TestImageMarshaler(new(feat.Gray))
	if err != nil {
		t.Error(err)
	}
}
