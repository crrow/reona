package reona

// TODO: GET RID OF IT, very disgusting

type bounder[K comparable] interface {
	flag() boundFlag
	key() K
}

type boundFlag uint8

const (
	included boundFlag = iota
	excluded
	unbounded
)

type baseBound[K comparable] struct {
	_key K
}

func (t baseBound[K]) key() K {
	return t._key
}

type includedBound[K comparable] struct {
	baseBound[K]
}

func (ib includedBound[K]) flag() boundFlag {
	return included
}

type excludedBound[K comparable] struct {
	baseBound[K]
}

func (ib excludedBound[K]) flag() boundFlag {
	return excluded
}

type unboundedBound[K comparable] struct{}

func (ib unboundedBound[K]) flag() boundFlag {
	return unbounded
}

func (ib unboundedBound[K]) key() K {
	panic("should not call me")
}

func newBound[K comparable](key K, flag boundFlag) bounder[K] {
	switch flag {
	case included:
		return includedBound[K]{
			baseBound[K]{
				_key: key,
			},
		}
	case excluded:
		return excludedBound[K]{
			baseBound[K]{
				_key: key,
			},
		}
	default:
		panic("should not call me")
	}
}

func newUnbound[K comparable]() bounder[K] {
	return unboundedBound[K]{}
}
