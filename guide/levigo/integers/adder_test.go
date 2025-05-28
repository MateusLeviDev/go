package integers

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(5, 5)
	expected := 10

	assertCorrectMessage(t, sum, expected)
}

func TestSplit(t *testing.T) {
	gotX, gotY := Split(9)
	expectedX, expectedY := 4, 5

	if gotX != expectedX || gotY != expectedY {
		t.Errorf("got %d and %d, expected %d and %d", gotX, gotY, expectedX, expectedY)
	}
}

func assertCorrectMessage(t testing.TB, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expected '%d' but got '%d'", expected, got)
	}
}
