package main

import (
	"fmt"
	"reflect"
)

// Type Assert
type Container []interface{}

func (c *Container) Put(elem interface{}) {
	*c = append(*c, elem)
}

func (c *Container) Get() interface{} {
	elem := (*c)[0]
	*c = (*c)[1:]
	return elem
}

// Reflection
type ReflectionContainer struct {
	s reflect.Value
}

func NewContainer(t reflect.Type, size int) *ReflectionContainer {
	if size <= 0 {
		size = 64
	}
	return &ReflectionContainer{s: reflect.MakeSlice(reflect.SliceOf(t), 0, size)}
}

func (c *ReflectionContainer) RPut(val interface{}) error {
	if reflect.ValueOf(val).Type() != c.s.Type().Elem() {
		return fmt.Errorf("Put: cannot put a %T into a slice of %s", val, c.s.Type().Elem())
	}
	c.s = reflect.Append(c.s, reflect.ValueOf(val))
	return nil
}
func (c *ReflectionContainer) RGet(refval interface{}) error {
	if reflect.ValueOf(refval).Kind() != reflect.Ptr ||
		reflect.ValueOf(refval).Elem().Type() != c.s.Type().Elem() {
		return fmt.Errorf("Get: needs *%s but got %T", c.s.Type().Elem(), refval)
	}
	reflect.ValueOf(refval).Elem().Set(c.s.Index(0))
	c.s = c.s.Slice(1, c.s.Len())
	return nil
}

func main() {
	intContainer := &Container{}
	intContainer.Put(7)
	intContainer.Put(42)
	elem, ok := intContainer.Get().(int)
	if !ok {
		fmt.Println("Unable to read an int from intContainer")
	}
	fmt.Printf("assertExample: %d (%T)\n", elem, elem)

	f1 := 3.1415926
	f2 := 1.41421356237
	c := NewContainer(reflect.TypeOf(f1), 16)
	if err := c.RPut(f1); err != nil {
		panic(err)
	}
	if err := c.RPut(f2); err != nil {
		panic(err)
	}
	g := 0.0
	if err := c.RGet(&g); err != nil {
		panic(err)
	}
	fmt.Printf("%v (%T)\n", g, g) //3.1415926 (float64)
	fmt.Println(c.s.Index(0))     //1.4142135623
}