package iter

type Iterator[T any] interface {
	Next() (T, bool)
}

type Error interface {
    Err() error
}

type Partitioner[T any] interface {
	Partition(int) []Iterator[T]
}

func Slice[T any](s []T) *slice[T] {
	return &slice[T]{data: s, currentEnd: len(s) - 1}
}

type slice[T any] struct {
	data       []T
	current    int
	currentEnd int
}

func (i *slice[T]) Next() (T, bool) {
	if i.current > i.currentEnd {
		var zero T
		return zero, false
	}
	i.current++
	return i.data[i.current-1], true
}

func (i *slice[T]) Partition(max int) []Iterator[T] {
	len := i.currentEnd - i.current + 1
	split := float64(len) / float64(max)
	idx := 0.0
	parts := []Iterator[T]{}

	for int(idx) < len {
		parts = append(parts, &slice[T]{data: i.data, current: int(idx) + i.current, currentEnd: int(idx+split-1) + i.current})
		idx += split
	}

	return parts
}

func Filter[T any, I Iterator[T]](it I, predicate func(T) (bool, error)) *filter[T, I] {
	return &filter[T, I]{parent: it, predicate: predicate}
}

type filter[T any, I Iterator[T]] struct {
	parent    I
	predicate func(T) (bool, error)
	err error
}

func (i *filter[T, I]) Next() (T, bool) {
	return i.iterate(i.parent.Next)
}

func (i *filter[T, I]) iterate(next func() (T, bool)) (T, bool) {
	for val, ok := next(); ok; val, ok = next() {
		if ok, err := i.predicate(val); ok {
			return val, ok
		} else if err != nil {
			i.err = err
			return val, false
		}
	}

	var zero T
	return zero, false
}

func (i *filter[T, I]) Err() error {
    if e, ok := interface{}(i.parent).(Error); ok {
        return e.Err()
    }

    return i.err
}

func Map[T, U any, I Iterator[T]](parent I, mapper func(T) (U, error)) *mapIt[T, U, I] {
	return &mapIt[T, U, I]{parent: parent, mapper: mapper}
}

type mapIt[T, U any, I Iterator[T]] struct {
	parent I
	mapper func(T) (U, error)
    err error
}

func (i *mapIt[T, U, I]) Next() (U, bool) {
	return i.iterate(i.parent.Next)
}

func (i *mapIt[T, U, I]) iterate(next func() (T, bool)) (U, bool) {
	val, ok := next()
	if !ok {
		var zero U
		return zero, false
	}

    mapped, err := i.mapper(val)
    if err != nil {
        i.err = err
        return mapped, false
    }

	return mapped, true
}

func (i *mapIt[T, U, I]) Err() error {
    if e, ok := interface{}(i.parent).(Error); ok {
        return e.Err()
    }

    return i.err
}

func FlatMap[T, U any, I Iterator[T], J Iterator[U]](parent I, mapper func(T) (J, error)) *flatMap[T, U, I, J] {
	return &flatMap[T, U, I, J]{parent: parent, mapper: mapper}
}

type flatMap[T, U any, I Iterator[T], J Iterator[U]] struct {
	parent I
	inner *J
	mapper func(T) (J, error)
	err error
}

func (i *flatMap[T, U, I, J]) Next() (U, bool) {
	if i.inner == nil {
		val, ok := i.parent.Next()
		if !ok {
			var zero U
			return zero, false
		}

		j, err := i.mapper(val)
		if err != nil {
			i.err = err
			var zero U
			return zero, false
		}

		i.inner = &j
	}

	val, ok := (*i.inner).Next()
	if !ok {
		i.inner = nil
		return i.Next()
	}

	return val, true
}

func (i *flatMap[T, U, I, J]) Err() error {
    if e, ok := interface{}(i.parent).(Error); ok {
        return e.Err()
    }

    return i.err
}

func ForEach[T any, I Iterator[T]](it I, consumer func(T)) error {
	for val, ok := it.Next(); ok; val, ok = it.Next() {
		consumer(val)
	}

    if e, ok := interface{}(it).(Error); ok {
        return e.Err()
    }

    return nil
}

func Reduce[T any, I Iterator[T]](it I, start T, accumulator func(T, T) T) (T, error) {
	result := start

	for val, ok := it.Next(); ok; val, ok = it.Next() {
		result = accumulator(result, val)
	}

    if e, ok := interface{}(it).(Error); ok {
        return result, e.Err()
    }

	return result, nil
}
