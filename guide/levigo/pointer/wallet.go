package main

import (
	"errors"
	"fmt"
	"sync"
)

type Bitcoin int

type Owner string

func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

type Wallet struct {
	mu      sync.Mutex
	owner   Owner
	balance Bitcoin
}

func NewWallet(owner Owner, balance Bitcoin) *Wallet {
	return &Wallet{owner: owner, balance: balance}
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if amount > w.balance {
		return ErrInsufficientFunds
	}

	w.balance -= amount
	return nil
}

func (w *Wallet) TransferTo(dest *Wallet, amount Bitcoin) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if amount > w.balance {
		return ErrInsufficientFunds
	}
	w.balance -= amount
	dest.balance += amount
	return nil
}
