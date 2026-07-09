package main

import (
	"errors"
	"fmt"
)

func main() {
	//x := 10
	/*	m := IncrementCopy(x)
		fmt.Println(m)
		fmt.Println(x)
	*/
	/*	origin := Wallet{Balance: 100}
		updated := WithDeposit(origin, 20)
		fmt.Println(origin.Balance)
		fmt.Println(updated.Balance)
		Deposit(&origin, 30)*/
	//if  {
	//
	//}

	//fmt.Println(origin.Balance)

	/*	a := NewCounter()
		b := NewCounter()
		fmt.Println(a())
		fmt.Println(a())
		fmt.Println(b())*/

	//double := MakeMultiplayer(2)
	//triple := MakeMultiplayer(3)
	//fmt.Println(double(10))
	//fmt.Println(triple(10))

}

func IncrementCopy(x int) int {
	x++
	return x
}

type Wallet struct {
	Balance int
}

func WithDeposit(wallet Wallet, amount int) Wallet {
	wallet.Balance += amount
	return wallet
}

func Deposit(wallet *Wallet, amount int) bool {
	if wallet == nil || amount <= 0 {
		return false
	}

	wallet.Balance += amount
	return true
}

type User struct {
	Name string
}

func Normalize(user *User) {

}

func NormalizeName(user *User) {

}
func NormalizedName(user User) string {
	return ""
}

func Deposit1(wallet *Wallet, amount int) error {
	if wallet == nil || amount <= 0 {
		return errors.New("error")
	}

	wallet.Balance += amount
	return nil
}

func Example() {
	status := "outer"
	{
		status := "inner"
		fmt.Println(status)
	}
	fmt.Println(status)
}

func CountPositive(values []int) int {
	var count int
	for _, value := range values {
		if value > 0 {
			// Ошибочно
			count := count + 1
			//_ = count
			// Корректно
			count = count + 1
		}
	}
	return count
}

var count int

func NewCounter() func() int {
	return func() int {
		count++
		return count
	}
}

func MakeMultiplayer(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}
