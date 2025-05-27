package main

import "testing"

func TestHallo(t *testing.T) {
	t.Run("saying Hallo to people", func(t *testing.T) {
		got := Hallo("levi")
		want := "Hallo, levi"

		assertCorrectMessage(t, got, want)
	})
	t.Run("say 'Hallo, Welt' when an empty string is supplied", func(t *testing.T) {
		got := Hallo("")
		want := "Hallo, Welt"

		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
