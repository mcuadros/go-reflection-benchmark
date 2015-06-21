package benchmark

import (
	"math/rand"
	"reflect"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	ints []int
}

var _ = Suite(&S{})

func (s *S) SetUpTest(c *C) {
	rand.Seed(42)

	s.ints = make([]int, 1000)
	for i := 0; i < 1000; i++ {
		s.ints[i] = rand.Intn(1e6)
	}
}

func (s *S) BenchmarkInitVar(c *C) {
	var value int
	for i := 0; i < c.N; i++ {
		value = 42
	}

	_ = value
}

func (s *S) BenchmarkInitVarReflect(c *C) {
	var value reflect.Value
	for i := 0; i < c.N; i++ {
		value = reflect.ValueOf(42)
	}

	_ = value
}

func (s *S) BenchmarkMakeSliceAndFill(c *C) {
	for i := 0; i < c.N; i++ {
		value := make([]int, 0, 0)
		for i := 0; i < 100; i++ {
			value = append(value, 42)
		}
	}
}

func (s *S) BenchmarkMakeSliceAndFillReflect(c *C) {
	for i := 0; i < c.N; i++ {
		value := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(42)), 0, 0)
		for i := 0; i < 100; i++ {
			value = reflect.Append(value, reflect.ValueOf(42))
		}
	}
}

func (s *S) BenchmarkMakeMapAndFill(c *C) {
	for i := 0; i < c.N; i++ {
		value := make(map[int]int, 0)
		for i := 0; i < 100; i++ {
			value[s.ints[i]] = 42
		}
	}
}

func (s *S) BenchmarkMakeMapAndFillReflect(c *C) {
	for i := 0; i < c.N; i++ {
		T := reflect.TypeOf(42)
		value := reflect.MakeMap(reflect.MapOf(T, T))
		for i := 0; i < 100; i++ {
			value.SetMapIndex(
				reflect.ValueOf(s.ints[i]),
				reflect.ValueOf(42),
			)
		}
	}
}

func (s *S) BenchmarkCallFunc(c *C) {
	for i := 0; i < c.N; i++ {
		multiply(42, 42)
	}
}

func (s *S) BenchmarkCallFuncReflect(c *C) {
	for i := 0; i < c.N; i++ {
		fn := reflect.ValueOf(multiply)
		fn.Call([]reflect.Value{reflect.ValueOf(42), reflect.ValueOf(42)})
	}
}

func (s *S) BenchmarkMakeChanAndPut(c *C) {
	for i := 0; i < c.N; i++ {
		ch := make(chan int)
		go func() {
			_ = <-ch
		}()

		ch <- 42
		close(ch)
	}
}

func (s *S) BenchmarkMakeChanAndPutReflect(c *C) {
	for i := 0; i < c.N; i++ {
		ch := reflect.MakeChan(reflect.TypeOf(make(chan int)), 0)
		go func() {
			ch.Recv()
		}()

		ch.Send(reflect.ValueOf(42))
		ch.Close()
	}
}

func multiply(a, b int) int {
	return a * b
}
