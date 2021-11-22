# iter
Experimental lazy iterator library to test the capabilities of go generics in 1.18

### Some takeaways

#### Generic methods with custom type parameters

The current implementation does not allow methods to specify different type parameters from the ones already specified in the type itself. Because of this restriction, the API design might skew towards functions, rather than methods, if some of them need different type parameters. This works around the limitation, and is not a showstopper, but it does lead to a bit more unweildy usage of the API, such as:

```go

s := Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

// If methods could define their own types:
s.Filter(func(i int) (bool, error) {
	return i%3 == 0, nil
}).Map(func(v int) (string, error) {
	return "number: " + strconv.Itoa(v), nil
})

// vs.

Map(Filter(s, func(i int) (bool, error) {
	return i%3 == 0, nil
}), func(v int) (string, error) {
	return "number: " + strconv.Itoa(v), nil
})
```

The key difference is that the code reads backwards, as one must start with the final operation, and drill down to the first one, whereas using methods allows the code to flow naturally from the start to the end operation.

Using temporary variables is always an option, and leads to more readability, perhaps in both cases, rather than just the current.

#### Generic methods with additional constraints on type parameters

Tangentially to the previous point, the design does not allow any changes to the constraints of the method's type parameters, compared to the type itself. Consider the following code if this was not the case:

```go
type BackwardsIterator[T any] interface {
	Iterator[T]
	Prev() (T, bool)
}

// an additional method to the filter type
func (i *filter[T any, I BackwardsIterator[T]]) Prev() (T, bool) {
	return i.iterate(i.parent.Prev)
}
```

This might allow concrete instances of this type to have the `Prev` method, only if their parent is also constrainted by the `BackwardsIterator` as well. For iteration sources like slices, this would allow the iteration to flow backwards, while for others, like channels, it would result in a compile me error. Currently, I can only think of a runtime solution to this, by checking if the parent supports the interface, which is not ideal.
