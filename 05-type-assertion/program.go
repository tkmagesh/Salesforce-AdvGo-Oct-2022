package main

import "fmt"

func main() {
	var x interface{}
	x = 100
	x = true
	x = "this is a string"
	x = 12.345
	x = struct{}{}
	fmt.Println(x)

	x = "this is a string"
	if val, ok := x.(int); !ok {
		fmt.Println("Not an int")
	} else {
		y := val + 100
		fmt.Println(y)
	}

	x = []int{}
	switch val := x.(type) {
	case int:
		fmt.Println("x is an int, x + 100 =", val+100)
	case string:
		fmt.Println("x is a string, len(x) = ", len(val))
	default:
		fmt.Println("unknown type")
	}

}
