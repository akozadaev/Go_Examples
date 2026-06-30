package main

import (
	"fmt"

	"example.com/functions/part2/pointervalue/account"
)

func main() {
	wallet := account.Wallet{Balance: 100}
	copyWithDeposit := account.WithDeposit(wallet, 50)

	fmt.Println("original after WithDeposit:", wallet.Balance)
	fmt.Println("copy after WithDeposit:", copyWithDeposit.Balance)

	account.Deposit(&wallet, 25)
	fmt.Println("original after Deposit:", wallet.Balance)

	a, b := 1, 2
	account.Swap(&a, &b)
	fmt.Println("swapped:", a, b)
}
