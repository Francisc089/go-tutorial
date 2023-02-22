// Declare a main package
// A package is a collection of source files that are compiled together
// A package is a way to group functions and it's made up of all the files in the same directory
package main

// Import the fmt package, which contains functions for formatting text, including printing to the console
// fmt is a standard library package
import (
	"fmt"

	"rsc.io/quote"
)

// Declare a main function. A main function executes by default when you run the program
func main() {
	fmt.Println(quote.Go())
}
