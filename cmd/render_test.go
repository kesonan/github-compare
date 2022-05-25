package cmd

import (
	"fmt"
	"testing"
	"time"
)

func TestRender(t *testing.T) {
	createAt, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-05-23 17:14:00", time.Local)
	// s := formatPeriod(time.Since(createAt))
	s := formatDuration(createAt)
	fmt.Println(s)
}
