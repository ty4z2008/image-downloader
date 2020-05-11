package utils

import (
	"testing"
)

func TestInit(t *testing.T) {
	err := Init("test1", false)
	if err != nil {
		t.Fatal("defaultLogger does not match logger returned by Init")
	}
}

func TestInfo(t *testing.T) {
	Info("test info ok")
}
