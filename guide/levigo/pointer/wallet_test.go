package main

import "testing"

func TestWallet(t *testing.T) {

	t.Run("deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("withdraw with funds", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("withdraw insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		wallet := Wallet{balance: startingBalance}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, startingBalance)
		assertError(t, err, ErrInsufficientFunds)
	})
}

func TestTransferTo(t *testing.T) {

	t.Run("transferTo with funds", func(t *testing.T) {
		origin := Wallet{Owner("Levi"), Bitcoin(100)}
		destiny := Wallet{owner: Owner("Maria"), balance: Bitcoin(50)}
		expectedBalance := Bitcoin(80)

		err := origin.TransferTo(&destiny, 30)

		assertBalance(t, destiny, expectedBalance)
		assertNoError(t, err)
	})

	t.Run("transferTo insufficient funds", func(t *testing.T) {
		startingBalance := Bitcoin(20)
		origin := Wallet{Owner("Levi"), Bitcoin(20)}
		destiny := Wallet{owner: Owner("Maria"), balance: Bitcoin(50)}

		err := origin.TransferTo(&destiny, 30)

		assertBalance(t, origin, startingBalance)
		assertError(t, err, ErrInsufficientFunds)
	})
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
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
