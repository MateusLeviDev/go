package main

import "testing"

func TestInventoryAdd(t *testing.T) {
	inventory := Inventory{"Product 1": 20}

	got, _ := inventory.Get("Product 1")
	want := 20

	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}
