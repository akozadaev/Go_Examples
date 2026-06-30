package account

import "testing"

func TestWithDepositReturnsCopy(t *testing.T) {
	original := Wallet{Balance: 100}
	updated := WithDeposit(original, 50)

	if original.Balance != 100 {
		t.Fatalf("original balance changed to %d", original.Balance)
	}
	if updated.Balance != 150 {
		t.Fatalf("updated balance = %d, want 150", updated.Balance)
	}
}

func TestWithDepositInvalidAmount(t *testing.T) {
	original := Wallet{Balance: 100}
	updated := WithDeposit(original, 0)
	if updated != original {
		t.Fatalf("WithDeposit() = %+v, want unchanged %+v", updated, original)
	}
}

func TestDepositMutatesWallet(t *testing.T) {
	wallet := Wallet{Balance: 100}
	if ok := Deposit(&wallet, 25); !ok {
		t.Fatal("Deposit() returned false, want true")
	}
	if wallet.Balance != 125 {
		t.Fatalf("wallet balance = %d, want 125", wallet.Balance)
	}
}

func TestDepositRejectsNilAndInvalidAmount(t *testing.T) {
	if ok := Deposit(nil, 10); ok {
		t.Fatal("Deposit(nil) returned true, want false")
	}

	wallet := Wallet{Balance: 100}
	if ok := Deposit(&wallet, -1); ok {
		t.Fatal("Deposit(invalid amount) returned true, want false")
	}
	if wallet.Balance != 100 {
		t.Fatalf("wallet balance = %d, want unchanged 100", wallet.Balance)
	}
}

func TestSwap(t *testing.T) {
	a, b := 1, 2
	if ok := Swap(&a, &b); !ok {
		t.Fatal("Swap() returned false, want true")
	}
	if a != 2 || b != 1 {
		t.Fatalf("Swap() got a=%d b=%d, want a=2 b=1", a, b)
	}
}

func TestSwapNil(t *testing.T) {
	a := 1
	if ok := Swap(&a, nil); ok {
		t.Fatal("Swap() returned true, want false")
	}
	if a != 1 {
		t.Fatalf("a = %d, want unchanged 1", a)
	}
}
