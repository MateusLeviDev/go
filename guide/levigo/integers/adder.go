package integers

func Add(x, y int) int {
	return x + y
}

func Split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}
