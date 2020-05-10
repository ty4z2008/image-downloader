package utils

import (
	"fmt"
	"strings"
)

func ByteCount(b int64, args ...string) string {
	var format string
	argStr := strings.Join(args, "")
	if argStr == "" {
		format = "B"
	} else {
		format = argStr
	}
	var unit int64 = 1
	switch format {
	case "KB":
		unit = 1024
	case "MB":
		unit = 1024 * 1024
	case "GB":
		unit = 1024 * 1024 * 1024
	default:
		unit = 1
	}

	return fmt.Sprintf("%.1f%s", float64(b)/float64(unit), format)
}
