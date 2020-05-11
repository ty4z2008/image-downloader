package utils

import (
	"testing"
)

func TestByteCount(t *testing.T) {

	t.Run("test convert to kb", func(t *testing.T) {
		res := ByteCount(4098, "KB")
		if res != "4.0KB" {
			t.Errorf("byteCount fail")
		}
	})
}
