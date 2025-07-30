package channel

func ProduceBy[Out any](fn func(out chan<- Out)) <-chan Out {
	out := make(chan Out)

	go func() {
		defer close(out)
		fn(out)
	}()

	return out
}

func TransformBy[In any, Out any](in <-chan In, fn func(in <-chan In, out chan<- Out)) <-chan Out {
	out := make(chan Out)

	go func() {
		defer close(out)
		fn(in, out)
	}()

	return out
}
