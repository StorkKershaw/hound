package toast

import (
	"strings"

	"github.com/electricbubble/go-toast"
)

type Writer struct{}

func (Writer) Write(p []byte) (int, error) {
	_ = toast.Push(
		strings.TrimRight(string(p), "\n"),
		toast.WithTitle("🐕 hound"),
		toast.WithAppID("hound"),
		toast.WithAudio(toast.Default),
	)
	return len(p), nil
}
