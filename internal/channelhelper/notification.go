package channelhelper

import (
	"github.com/electricbubble/go-toast"
)

func Notify(in <-chan error, out chan<- struct{}) {
	toast.Push(
		"🐕ᯓ hound is running",
		toast.WithTitle("🐕 hound"),
		toast.WithAppID("hound"),
		toast.WithAudio(toast.Default),
	)

	for err := range in {
		_ = toast.Push(
			err.Error(),
			toast.WithTitle("🐕 hound"),
			toast.WithAppID("hound"),
			toast.WithAudio(toast.Default),
		)
	}

	out <- struct{}{}
}
