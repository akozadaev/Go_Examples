package main

import (
	"fmt"
	"unsafe"
)

type Base struct {
	value int
}

func (b *Base) Print() {
	fmt.Println("Base: ", b.value)
}

type Derived struct {
	Base
}

func (d *Derived) Print() {
	fmt.Println("Base: ", d.value)
}

type BadStruct struct {
	a bool  // 1
	b int64 //8
	c bool  // 16
}
type GoodStruct struct {
	b int64
	a,
	c bool
}

func main() {

	var bs BadStruct
	var gs GoodStruct
	fmt.Println("Size bs: ", unsafe.Sizeof(bs))
	fmt.Println("Size gs: ", unsafe.Sizeof(gs))
	d := Derived{Base{value: 3223}}
	d.Print()
	d.Base.Print()

	/*counter := Counter{0}
	counter.Increment()
	fmt.Println(counter.count)*/
	/*
		acc := account.New(9848745)
		fmt.Println(acc.Balance())
		if err := acc.Deposit(500); err != nil {
			log.Fatal(err)
		}
		fmt.Println(acc.Balance())
		if err := acc.Withdraw(984874500); err != nil {
			log.Fatal(err)
		}*/
	/*	type Animal struct {
			Name string
		}
		type Dog struct {
			Animal
		}*/

	/*	//var p Point
		// Литерал
		p := Point{1.0, 2.2}
		// Литерал + имя поля
		p1 := Point{X: 1.0}
		p2 := Point{Y: 1.0, X: 2.2}
		p3 := NewPoint(1.0, 1.0)
		fmt.Println(p, p1, p2, p3)*/
	/*	c1 := Color{255, 0, 0}
		c2 := Color{255, 0, 0}
		c3 := Color{0, 0, 255}
		fmt.Println(c1 == c2)
		fmt.Println(c1 == c3)*/

	/*	p := Point{X: 1.0}
		ptr := &p
		ptr.X = 1212
		fmt.Println(ptr.X)
		(*ptr).X = 1213
		fmt.Println(ptr.X)

		modify1(p)
		fmt.Println(p)
		modify2(ptr)
		fmt.Println(ptr)*/

	//config := struct {
	//	Host string
	//	Port int
	//}{
	//	Host: "localhost",
	//	Port: 8080,
	//}
	//fmt.Println(config)
	//u := account.Account{Name: "dfdfd"}

}

/*
type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

type Color struct {
	R, G, B uint8
}

// не влияет на оригинал
func modify1(p Point) {
	p.X = 898
}

// влияет на оригинал
func modify2(p *Point) {
	p.X = 654567
}
*/
