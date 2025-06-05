package main

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})
}

func TestSumAll(t *testing.T) {
	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func checkSums(t testing.TB, got, want []int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSumAllTails(t *testing.T) {
	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("safely sum slices with empty, single and multiple items", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{1}, []int{10, 20, 30})
		want := []int{0, 0, 50}
		checkSums(t, got, want)
	})
}

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }

	t.Run("safely filter", func(t *testing.T) {
		got := Filter([]int{1, 2, 3, 4, 5, 6}, isEven)
		want := []int{2, 4, 6}
		checkSums(t, got, want)
	})

	t.Run("empty slice", func(t *testing.T) {
		got := Filter([]int{}, isEven)
		want := []int{}
		checkSums(t, got, want)
	})
}
