package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

func main() {
	x := 324
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	fmt.Println(t)
	fmt.Println(v.Int())
	// === Проблема typed nil ====
	var a any = nil
	var p *int = nil
	var b any = p
	fmt.Println("reflect.TypeOf(a) =", reflect.TypeOf(a))
	fmt.Println("reflect.ValueOf(a) =", reflect.ValueOf(a))
	fmt.Println("reflect.TypeOf(b) =", reflect.TypeOf(b))
	fmt.Println("reflect.ValueOf(b) =", reflect.ValueOf(b))
	fmt.Println(a, b)

	err := f()
	fmt.Println(err)

	// Тип, имя, категория
	fmt.Println("Тип, имя, категория")

	id := UserID(7)
	tp := reflect.TypeOf(id)

	fmt.Println(tp.String()) // main.UserID
	fmt.Println(tp.Name())   // UserID
	fmt.Println(tp.PkgPath())
	fmt.Println(tp.Kind()) // int64

	describe(id)
	// TypeFor
	writerType := reflect.TypeFor[io.Writer]()
	fileType := reflect.TypeFor[*os.File]()
	fmt.Println(fileType.Implements(writerType)) // true

	// AssignableTo
	type UserID int64

	userIDType := reflect.TypeFor[UserID]()
	int64Type := reflect.TypeFor[int64]()

	fmt.Println(userIDType.AssignableTo(int64Type))  // false
	fmt.Println(userIDType.ConvertibleTo(int64Type)) // true

	// Поля структуры
	type PersonExample struct {
		Name  string
		Age   int
		Email string
	}
	pe := PersonExample{Name: "Alice", Age: 30, Email: "alice@example.com"}
	ts := reflect.TypeOf(pe)
	//До Go 1.26 поля структуры обычно перебирали индексами:
	for i := 0; i < ts.NumField(); i++ {
		field := ts.Field(i)
		fmt.Println(field.Name, field.Type)
	}

	//В Go 1.26 появились методы-итераторы Type.Fields, Type.Methods, Type.Ins и Type.Outs. Они сочетаются с range над итераторами:
	for field := range ts.Fields() {
		fmt.Println(field.Name, field.Type)
	}

	// reflect.Value: чтение значения
	x1 := 42
	v1 := reflect.ValueOf(x1)

	fmt.Println(v1.Type())      // int
	fmt.Println(v1.Kind())      // int
	fmt.Println(v1.Int())       // 42
	fmt.Println(v1.Interface()) // 42 как значение типа any

	//Валидность
	//Нулевой reflect.Value ничего не представляет:
	var v2 reflect.Value
	fmt.Println(v2.IsValid()) // false

	fieldByName(pe)

	//	Адресуемость и изменение значений
	x3 := 42
	v3 := reflect.ValueOf(&x3)
	elem := v3.Elem()

	fmt.Println(elem.CanAddr()) // true
	fmt.Println(elem.CanSet())  // true
	elem.SetInt(24)
	fmt.Println(x3) // 24

	//	Создание значений во время выполнения
	type User struct {
		Name string
		Age  int
	}

	userType := reflect.TypeFor[User]()
	ptr := reflect.New(userType) // *User
	user := ptr.Elem()           // User

	user.FieldByName("Name").SetString("Анна")
	user.FieldByName("Age").SetInt(30)

	created := ptr.Interface().(*User)
	fmt.Println(created)

	sliceType := reflect.SliceOf(reflect.TypeFor[string]())
	slice := reflect.MakeSlice(sliceType, 0, 4)
	slice = reflect.Append(slice, reflect.ValueOf("Go"))
	slice = reflect.Append(slice, reflect.ValueOf("1.26"))

	//	Структуры и теги
	type UserStr struct {
		ID       int    `json:"id" db:"user_id"`
		Name     string `json:"name,omitempty"`
		Password string `json:"-"`
		token    string
	}

	// Создаём экземпляр структуры
	u := UserStr{
		ID:       123,
		Name:     "Alice",
		Password: "secret",
		token:    "internal",
	}

	// Вызываем toMap с экземпляром, а не с типом
	resultMap, errMap := toMap(u)
	if errMap != nil {
		fmt.Println("Ошибка toMap:", errMap)
	} else {
		fmt.Println("Результат toMap:", resultMap)
	}

	//	Функции и методы

}

type UserID int64
type MyError struct{}

func (*MyError) Error() string { return "ошибка" }

func f() error {
	var err *MyError
	return err
}

func describe(x any) {
	t := reflect.TypeOf(x)
	if t == nil {
		fmt.Println("nil interface")
		return
	}

	fmt.Printf("type=%v kind=%v\n", t, t.Kind())

	switch t.Kind() {
	case reflect.Struct:
		fmt.Println("fields:", t.NumField())
	case reflect.Slice, reflect.Array:
		fmt.Println("element type:", t.Elem())
	case reflect.Map:
		fmt.Println("key:", t.Key(), "value:", t.Elem())
	case reflect.Func:
		fmt.Println("inputs:", t.NumIn(), "outputs:", t.NumOut())
	case reflect.Int64:
		fmt.Println("basic type of kind", t.Kind())
	default:
		panic("unhandled default case")
	}
}

func fieldByName(v any) {
	val := reflect.ValueOf(v)
	field := val.FieldByName("Age")
	if !field.IsValid() {
		fmt.Println("field Age not found")
		return
	}
	fmt.Println("field Age value:", field.Interface())
}

func toMap(input any) (map[string]any, error) {
	v := reflect.ValueOf(input)
	if !v.IsValid() {
		return nil, errors.New("nil input")
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, errors.New("nil pointer")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("want struct, got %s", v.Kind())
	}

	result := make(map[string]any)
	for sf, fv := range v.Fields() {
		if !sf.IsExported() || !fv.CanInterface() {
			continue
		}
		name, options, _ := strings.Cut(sf.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = sf.Name
		}
		if options == "omitempty" && fv.IsZero() {
			continue
		}
		result[name] = fv.Interface()
	}
	return result, nil
}
