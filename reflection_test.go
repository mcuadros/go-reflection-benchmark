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
	for i := 0; i < c.N; i++ {
		reflect.ValueOf(42)
	}

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
		MakeSliceAndFill()
	}
}

func (s *S) BenchmarkMakeMapAndFill(c *C) {
	for i := 0; i < c.N; i++ {
		value := make(map[int]int, 0)
		for i := 0; i < 100; i++ {
			value[i] = 42
		}
	}
}

func (s *S) BenchmarkMakeMapAndFillReflect(c *C) {
	for i := 0; i < c.N; i++ {
		MakeMapAndFill()
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
		GetMapKeys(m)
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
	impl := func(args []reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf(
			int(args[0].Int()) * int(args[1].Int()),
		)}
	}

	for i := 0; i < c.N; i++ {
		MakeFuncAndCall(impl)
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
		MakeChanAndPut()
	}
}

func (s *S) BenchmarkNewStructAndSetFieldValue(c *C) {
	for i := 0; i < c.N; i++ {
		s := new(Foo)
		s.Value = 42
	}
}

func (s *S) BenchmarkNewStructAndSetFieldValueReflect(c *C) {
	for i := 0; i < c.N; i++ {
		NewStructAndSetFieldValue()
	}
}

func (s *S) BenchmarkCallStructMethod(c *C) {
	value := &Foo{42}
	for i := 0; i < c.N; i++ {
		value.Multiply(42)
	}
}

func (s *S) BenchmarkCallStructMethodReflect(c *C) {
	value := NewStructAndSetFieldValue()
	for i := 0; i < c.N; i++ {
		CallStructMethod(value)
	}
}

func (s *S) buildMap(l int) map[int]int {
	value := make(map[int]int, 0)
	for i := 0; i < l; i++ {
		value[s.ints[i]] = 42
	}

	return value
}
