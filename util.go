package main

import (
	"fmt"
	"os"
)

// Read int from a string and exit on fail
func SScanInt(s string, i *int, name string) {
	if scanned, err := fmt.Sscanf(s, "%d", i); err != nil || scanned < 1 {
		fmt.Println("Failed to read param")
		os.Exit(1)
	}

	if *i <= 0 {
		fmt.Printf("%s must be a positive integer\n", name)
		os.Exit(1)
	}
}

// Read int from the user and exit on fail
func ScanInt(i *int, name string) {
	if scanned, err := fmt.Scanf("%d", i); err != nil || scanned < 1 {
		fmt.Println("Failed to read param")
		os.Exit(1)
	}

	if *i <= 0 {
		fmt.Printf("%s must be a positive integer\n", name)
		os.Exit(1)
	}
}

// Read int from a string and exit on fail
func SScanFloat(s string, f *float64, name string) {
	if scanned, err := fmt.Sscanf(s, "%f", f); err != nil || scanned < 1 {
		fmt.Println("Failed to read param")
		os.Exit(1)
	}

	if *f <= 0 {
		fmt.Printf("%s must be a positive number\n", name)
		os.Exit(1)
	}
}

// Read float from the user and exit on fail
func ScanFloat(f *float64, name string) {
	if scanned, err := fmt.Scanf("%f", f); err != nil || scanned < 1 {
		fmt.Println("Failed to read param")
		os.Exit(1)
	}

	if *f <= 0 {
		fmt.Printf("%s must be a positive number\n", name)
		os.Exit(1)
	}
}
