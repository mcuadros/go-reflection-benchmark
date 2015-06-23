package benchmark

import "reflect"

func MakeSliceAndFill() {
	value := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(42)), 0, 0)
	for i := 0; i < 100; i++ {
		value = reflect.Append(value, reflect.ValueOf(42))
	}
}

func MakeMapAndFill() {
	T := reflect.TypeOf(42)
	value := reflect.MakeMap(reflect.MapOf(T, T))
	for i := 0; i < 100; i++ {
		value.SetMapIndex(
			reflect.ValueOf(i),
			reflect.ValueOf(42),
		)
	}
}

func GetMapKeys(m map[int]int) {
	value := reflect.ValueOf(m)
	value.MapKeys()
}

// MakeFuncAndCall
// `reflect.MakeFunc` creates a new function of a given signature, the returned
// function is a wrapper of other function with the following fixed signature:
//   func(args []reflect.Value) []reflect.Value
//
// a example of a implementation for simulate the function `multiply` is:
//   func(args []reflect.Value) []reflect.Value {
//        return []reflect.Value{reflect.ValueOf(
//              int(args[0].Int()) * int(args[1].Int()),
//        )}
//    }
func MakeFuncAndCall(i func(args []reflect.Value) []reflect.Value) {
	fn := reflect.MakeFunc(reflect.TypeOf(multiply), i)
	fn.Call([]reflect.Value{reflect.ValueOf(42), reflect.ValueOf(42)})
}

func MakeChanAndPut() {
	ch := reflect.MakeChan(reflect.TypeOf(make(chan int)), 0)
	go func() {
		ch.Recv()
	}()

	ch.Send(reflect.ValueOf(42))
	ch.Close()
}

func NewStructAndSetFieldValue() reflect.Value {
	value := reflect.New(reflect.TypeOf(Foo{}))

	f := value.Elem().FieldByName("Value")
	f.SetInt(42)

	return value
}

func CallStructMethod(value reflect.Value) {
	m := value.MethodByName("Multiply")
	m.Call([]reflect.Value{reflect.ValueOf(42)})
}

func multiply(a, b int) int {
	return a * b
}

type Foo struct {
	Value int
}

func (f *Foo) Multiply(factor int) int {
	return f.Value * factor
}
