package pages

import (
	"testing"
	"time"
)

const timeLayout = time.RFC3339

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	v, err := time.Parse(timeLayout, s)
	if err != nil {
		t.Fatal(err)
	}
	return v
}
