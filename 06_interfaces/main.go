package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

// 1. Утиная типизация.
type Speaker interface {
	Speak() string
}

type Dog struct{}

type Cat struct{}

type Robot struct{}

func (Dog) Speak() string   { return "woof" }
func (Cat) Speak() string   { return "meow" }
func (Robot) Speak() string { return "beep" }

func demoDuckTyping() {
	fmt.Println("=== 1. Утиная типизация ===")
	creatures := []Speaker{Dog{}, Cat{}, Robot{}}
	for _, c := range creatures {
		fmt.Println(c.Speak())
	}
	fmt.Println()
}

// 2. Композиция интерфейсов.
type Reader interface {
	Read([]byte) (int, error)
}

type Writer interface {
	Write([]byte) (int, error)
}

type ReadWriter interface {
	Reader
	Writer
}

type Buffer struct {
	data []byte
}

func (b *Buffer) Read(p []byte) (int, error) {
	if len(b.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func demoComposition() {
	fmt.Println("=== 2. Композиция интерфейсов ===")
	var rw ReadWriter = &Buffer{}
	_, _ = rw.Write([]byte("hello"))
	buf := make([]byte, 8)
	n, err := rw.Read(buf)
	if err != nil {
		fmt.Println("read error:", err)
		return
	}
	fmt.Println("read bytes:", string(buf[:n]))
	fmt.Println()
}

// 3. Утверждение типов и type switch.
func describeValue(v any) {
	switch x := v.(type) {
	case nil:
		fmt.Println("value is nil")
	case int:
		fmt.Printf("int: %d\n", x)
	case string:
		fmt.Printf("string: %q\n", x)
	case bool:
		fmt.Printf("bool: %v\n", x)
	default:
		fmt.Printf("unknown type: %T\n", x)
	}
}

func safeAssertion(v any) {
	if s, ok := v.(string); ok {
		fmt.Println("safe assertion -> string:", s)
	} else {
		fmt.Println("safe assertion -> not a string")
	}
}

func demoTypeAssertions() {
	fmt.Println("=== 3. Утверждение типов и type switch ===")
	describeValue(42)
	describeValue("hello")
	describeValue(true)
	describeValue(struct{}{})
	safeAssertion("go")
	safeAssertion(7)
	fmt.Println()
}

// 4. Внутреннее устройство интерфейса и nil trap.
type MyError struct {
	Msg string
}

func (e *MyError) Error() string { return e.Msg }

func demoNilInterfaces() {
	fmt.Println("=== 4. Внутреннее устройство интерфейса и nil trap ===")
	var i any = nil
	fmt.Println("i == nil:", i == nil)

	var p *int
	var j any = p
	fmt.Println("j == nil:", j == nil)
	fmt.Printf("j has dynamic type %T\n", j)

	var err error = (*MyError)(nil)
	fmt.Println("err == nil:", err == nil)
	fmt.Printf("err type: %T\n", err)
	fmt.Println()
}

// 5. Правила присваивания и ошибки.
var ErrNotFound = errors.New("not found")

func demoAssignmentsAndErrors() {
	fmt.Println("=== 5. Правила присваивания и ошибки ===")
	var w io.Writer
	w = os.Stdout
	fmt.Println("assigned writer ok:", w != nil)

	var s Speaker
	s = Dog{}
	fmt.Println("assigned speaker ok:", s.Speak())

	wrapped := fmt.Errorf("read failed: %w", ErrNotFound)
	fmt.Println("errors.Is:", errors.Is(wrapped, ErrNotFound))
	fmt.Println()
}

// 6. Безопасное и опасное приведение типов.
func demoCasting() {
	fmt.Println("=== 6. Безопасное и опасное приведение типов ===")
	var v any = "hello"
	if s, ok := v.(string); ok {
		fmt.Println("safe cast ->", s)
	}

	// Опасное приведение — может вызвать panic.
	// fmt.Println(v.(int))

	buf := bytes.NewBufferString("demo")
	var writer io.Writer = buf
	if b, ok := writer.(*bytes.Buffer); ok {
		fmt.Println("buffer type confirmed:", b.String())
	}
	fmt.Println()
}

// 7. Производительность и практические рекомендации.
func demoPerformanceNote() {
	fmt.Println("=== 7. Производительность и практические рекомендации ===")
	fmt.Println("Интерфейсы добавляют косвенность, поэтому в узких циклах стоит измерять производительность.")
	fmt.Println("Если важны гибкость и расширяемость, интерфейсы обычно оправданы.")
	fmt.Println()
}

func main() {
	demoDuckTyping()
	demoComposition()
	demoTypeAssertions()
	demoNilInterfaces()
	demoAssignmentsAndErrors()
	demoCasting()
	demoPerformanceNote()
}
