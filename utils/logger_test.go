package utils

import (
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("project name", false, func(t *testing.T) {
		Info("ok")
	})
}
