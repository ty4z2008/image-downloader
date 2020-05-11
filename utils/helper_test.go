package utils

import (
	"math/rand"
	"testing"
	"time"
)

func TestByteCount(t *testing.T) {

	t.Run("test convert to kb", func(t *testing.T) {
		res := ByteCount(4098, "KB")

	})
}
