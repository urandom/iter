package iter

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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

func work(v int) {
	time.Sleep(time.Millisecond)
}

func BenchmarkForEach(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rng := Range[int](0, 1_000, 1)
		b.StartTimer()

		ForEach(rng, work)
	}
}

func BenchmarkStream0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rng := Range[int](0, 1_000, 1)
		b.StartTimer()

		c := Stream[int](rng, 0)
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				for v := range c {
					work(v.Value)
				}
				wg.Done()
			}()
		}
	}
}
