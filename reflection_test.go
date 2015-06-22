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

func (s *S) BenchmarkGetMapKeys(c *C) {
	m := s.buildMap(100)
	for i := 0; i < c.N; i++ {
		keys := make([]int, 100)
		for k, _ := range m {
			keys = append(keys, k)
		}
	}
}

func (s *S) BenchmarkGetMapKeysReflect(c *C) {
	m := s.buildMap(100)
	for i := 0; i < c.N; i++ {
		reflect.ValueOf(m).MapKeys()
	}
}

func (s *S) BenchmarkCallFunc(c *C) {
	for i := 0; i < c.N; i++ {
		multiply(42, 42)
	}
}

func (s *S) BenchmarkMakeFuncAndCall(c *C) {
	base := func(args []int) int { return multiply(args[0], args[1]) }
	for i := 0; i < c.N; i++ {
		fn := func(a int, b int) int {
			return base([]int{a, b})
		}

		fn(42, 42)
	}
}

func (s *S) BenchmarkMakeFuncAndCallReflect(c *C) {
	base := func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(
			multiply(int(args[0].Int()), int(args[1].Int())),
		)}
	}

	for i := 0; i < c.N; i++ {
		fn := reflect.MakeFunc(reflect.TypeOf(multiply), base)
		fn.Call([]reflect.Value{reflect.ValueOf(42), reflect.ValueOf(42)})
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

func (s *S) buildMap(l int) map[int]int {
	value := make(map[int]int, 0)
	for i := 0; i < l; i++ {
		value[s.ints[i]] = 42
	}

	return value
}

func multiply(a, b int) int {
	return a * b
}
