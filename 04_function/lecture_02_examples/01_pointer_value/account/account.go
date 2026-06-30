package account

// Wallet хранит баланс в целых единицах для небольшого примера с указателями и значениями.
type Wallet struct {
	Balance int
}

// WithDeposit возвращает копию wallet с добавленной суммой amount.
func WithDeposit(wallet Wallet, amount int) Wallet {
	if amount <= 0 {
		return wallet
	}
	wallet.Balance += amount
	return wallet
}

// Deposit добавляет amount к wallet на месте. Возвращает false для nil или некорректной суммы.
func Deposit(wallet *Wallet, amount int) bool {
	if wallet == nil || amount <= 0 {
		return false
	}
	wallet.Balance += amount
	return true
}

// Swap меняет местами два значения int на месте. Возвращает false, если любой указатель равен nil.
func Swap(a, b *int) bool {
	if a == nil || b == nil {
		return false
	}
	*a, *b = *b, *a
	return true
}
