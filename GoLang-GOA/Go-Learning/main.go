package main

import "fmt"

var version string = "1.0.0"

const (
	author  = "Go Developer"
	license = "MIT"
)

func main() {
	/*fmt.Println("Hello, World!")

	var name string = "GoLang"
	fmt.Printf("Welcome to %s programming!\n", name)
	fmt.Printf("Version: %s, Author: %s, License: %s\n", version, author, license)

	// Demonstrating variable shadowing
	{
		var name string = "Gopher"
		fmt.Printf("Inside block, name is: %s\n", name)
	}
	fmt.Printf("Outside block, name is: %s\n", name)


	// Using short variable declaration
	count := 10
	fmt.Printf("Count is: %d\n", count)

	// array initialization
	numbers := [5]int{1, 2, 3, 4, 5}
	fmt.Printf("Numbers: %v\n", numbers)

	days := [...]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	fmt.Printf("Days of the week: %v\n", days)

	// iterating over an array
	for i, day := range days {
		fmt.Printf("Day %d: %s\n", i+1, day)
	}

	for i := 0; i < len(numbers); i++ {
		fmt.Printf("Number %d: %d\n", i+1, numbers[i])
	}

	// sum of array elements
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	fmt.Printf("Sum of numbers: %d\n", sum)


	// Demonstrating error handling
	result, err := divide(10, 3)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Result: %f\n", result)
	}

	*/

	var z int = 10
	fmt.Printf("Value of z: %d\n", z)
	x, y := 10, 10
	// comparison operators
	fmt.Printf("x == y: %t\n", x == y)
	fmt.Printf("x != y: %t\n", x != y)
	fmt.Printf("x > y: %t\n", x > y)
	fmt.Printf("x == z: %t\n", x == z)

}

func divide(a, b float32) (float32, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero is not allowed")
	}
	return a / b, nil
}
