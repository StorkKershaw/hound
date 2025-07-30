package main

import (
	"github.com/electricbubble/go-toast"
)

func notify(in <-chan error, out chan<- struct{}) {
	for err := range in {
		_ = toast.Push(
			err.Error(),
			toast.WithTitle("🐕 hound"),
			toast.WithAudio(toast.Default),
		)
	}

	out <- struct{}{}
}
