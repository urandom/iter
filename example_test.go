package iter

import (
	"fmt"
	"strconv"
	"strings"
)

func ExampleUsage() {
	s := Slice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	f := Filter(
		s,
		func(i int) (bool, error) {
			return i%3 == 0, nil
		},
	)

	m := Map(f, func(v int) (string, error) {
		return "number: " + strconv.Itoa(v), nil
	})

	flat := FlatMap[string, string, Iterator[string], Iterator[string]](m, func(v string) (Iterator[string], error) {
		return Slice(strings.Split(v, "")), nil
	})

	ForEach(flat, func(v string) {
		fmt.Printf("%v ", v)
	})

	res, err := Reduce(Slice([]int{1,2,3,4,5,6,7,8,9,10}), 0, func(acc, v int) int {
		return acc + v
	})

	fmt.Println("Reduce error:", err == nil)
	fmt.Println("Reduce sum:", res)

	res, _ = Reduce(Range[int](10, 100, 5), 0, func(acc, v int) int {
		return acc + v
	})

	fmt.Println("Reduce range:", res)

	// Output: n u m b e r :   3 n u m b e r :   6 n u m b e r :   9 Reduce error: true
	// Reduce sum: 55
	// Reduce range: 1045
}
