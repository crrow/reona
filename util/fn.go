package util

type Option[T any] interface {
	apply(opts *T)
}

type OptionFunc[T any] func(*T)

func (fn OptionFunc[T]) apply(opts *T) {
	fn(opts)
}

func ApplyOptions[T any](t *T, opts ...Option[T]) {
	for i := range opts {
		opts[i].apply(t)
	}
}
