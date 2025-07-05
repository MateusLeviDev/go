package main

import (
	"sync"
	"testing"
)

func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := NewWallet("Levi", 0)
		wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("test safely concurrent deposit", func(t *testing.T) {
		wantedCount := 1000
		wallet := NewWallet("Levi", 0)
		var wg sync.WaitGroup
		wg.Add(wantedCount)

		for i := 0; i < wantedCount; i++ {
			go func() {
				wallet.Deposit(Bitcoin(1))
				wg.Done()
			}()
		}
		wg.Wait()

		assertBalance(t, wallet, Bitcoin(1000))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := NewWallet("Levi", 20)
		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		wallet := NewWallet("Levi", 20)
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, Bitcoin(20))
		assertError(t, err, ErrInsufficientFunds)
	})
}

func TestTransferTo(t *testing.T) {

	t.Run("transferTo with funds", func(t *testing.T) {
		//origin := Wallet{Owner("Levi"), Bitcoin(100)}
		origin := NewWallet(Owner("Levi"), Bitcoin(100))
		destiny := NewWallet(Owner("Maria"), Bitcoin(50))
		expectedBalance := Bitcoin(80)

		err := origin.TransferTo(destiny, 30)

		assertBalance(t, destiny, expectedBalance)
		assertNoError(t, err)
	})

	t.Run("transferTo insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		origin := NewWallet(Owner("Levi"), Bitcoin(20))
		destiny := NewWallet(Owner("Maria"), Bitcoin(50))

		err := origin.TransferTo(destiny, 30)

		assertBalance(t, origin, startingBalance)
		assertError(t, err, ErrInsufficientFunds)
	})
}

func TestBitcoinStringer(t *testing.T) {
	b := Bitcoin(42)
	got := b.String()
	want := "42 BTC"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertBalance(t testing.TB, wallet *Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
