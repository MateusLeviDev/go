package main

func twoSumOQuadratic(nums []int, target int) []int {
	//O(n2) time complex
	//O(1) space complex
	for i, num1 := range nums {
		for j, num2 := range nums {
			if i != j {
				if num1+num2 == target {
					return []int{i, j}
				}
			}
		}
	}

	return nil
}

func twoSum(nums []int, target int) []int {
	pairs := make(map[int]int)

	for i, num := range nums {
		if val, ok := pairs[num]; ok {
			return []int{val, i} 
		} else {
			pairs[target-num] = i
		}
	}

	return nil
}

func main() {
	list := []int{2, 8, 7, 15}
	result := twoSum(list, 10)
	println(result[0], result[1])
}
