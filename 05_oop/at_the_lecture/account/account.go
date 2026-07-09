package account

import "errors"

type Account struct {
	balance float64
}

var ErrEmptyName = errors.New("error")
var ErrNonPositiveDeposit = errors.New("deposit amount must be positive")
var ErrNonPositiveWithdraw = errors.New("withdraw amount must be positive")
var ErrNonInsufficient = errors.New("insufficient")

func New(initBalance float64) *Account {
	if initBalance < 0 {
		initBalance = 0
	}
	return &Account{balance: initBalance}
}

func (a *Account) Balance() float64 {
	return a.balance
}

// Setter с валидацией
func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrNonPositiveDeposit
	}
	a.balance += amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrNonPositiveWithdraw
	}
	if a.balance < amount {
		return ErrNonInsufficient
	}
	a.balance -= amount
	return nil
}
