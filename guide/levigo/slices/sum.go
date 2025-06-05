package main

func Sum(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum += number
	}
	return sum
}

func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}
	return sums
}

func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}
	return sums
}

func Filter(numbers []int, fn func(int) bool) []int {
	seen := make(map[int]bool)
	result := make([]int, 0)
	for _, number := range numbers {
		if fn(number) && !seen[number] {
			seen[number] = true
			result = append(result, number)
		}
	}
	return result
}
