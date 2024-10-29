// Challenge algorithm implementations in Golang
// Author: Tommi Korpelainen

// This is the main file

package main

import (
	"fmt"
	"os"
)

// Program entry
func main() {

	var P int

	// Extract args
	for i, v := range os.Args {
		switch v {
		case "-P", "--P":
			SScanInt(os.Args[i+1], &P, "P")
		}
	}

	fmt.Println("==== Algorithm demo ====")

	if P == 0 {
		fmt.Println("Choose the program to execute:")
		fmt.Println("1. Graph coloring")
		fmt.Println("2. Local minimum")
		fmt.Println("3. Largest contiguous submatrix")
		fmt.Print("Program: ")
		ScanInt(&P, "Program")
	}

	// Choose the problem to demo
	switch P {
	case 1:
		RunGraphColor()
	case 2:
		RunLocalMinimum()
	case 3:
		RunSubmatrix()
	default:
		fmt.Println("Not a recognized program")
		os.Exit(1)
	}

}
