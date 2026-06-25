package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x int = 42
	var y interface{} = "hello"
	var err error = fmt.Errorf("oops")

	fmt.Println(reflect.TypeOf(x))   // int
	fmt.Println(reflect.TypeOf(y))   // string (фактический тип значения в интерфейсе)
	fmt.Println(reflect.TypeOf(err)) // *errors.errorString

	// Можно получить имя типа:
	t := reflect.TypeOf(x)
	fmt.Println(t.Name()) // int
	fmt.Println(t.Kind()) // int (категория типа)

	// Для пользовательских типов:
	type MyInt int
	var z MyInt = 10
	fmt.Println(reflect.TypeOf(z).Name()) // MyInt
	fmt.Println(reflect.TypeOf(z).Kind()) // int
}
