package featset_test

import (
	"testing"

	"github.com/jvlmdr/go-cv/featset"
)

func TestImageMarshaler(t *testing.T) {
	var err error
	err = featset.TestImageMarshaler(new(featset.RGB))
	if err != nil {
		t.Error(err)
	}
	err = featset.TestImageMarshaler(new(featset.Gray))
	if err != nil {
		t.Error(err)
	}
}
