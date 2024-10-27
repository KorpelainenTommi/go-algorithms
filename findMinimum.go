package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

// Find the index of the local minimum i, where f[i-1] >= f[i] <= f[i+1]
func FindLocalMinimum(f []float64, noPrint bool) int {
	start, end := 0, len(f)

	for {
		i := (end-start)/2 + start
		if !noPrint {
			fmt.Println("start:", start, "end:", end, "i:", i)
		}
		a, b := f[i-1], f[i+1]
		switch {
		case a >= f[i] && f[i] <= b:
			return i
		case b < f[i]:
			start = i
		default:
			end = i
		}
	}
}

// Save generated range to file
func SaveRangeJson(f []float64, path string) error {

	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to open", path)
		return err
	}

	data, err := json.Marshal(f)
	if err != nil {
		fmt.Println("Failed to json encode")
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Failed to write")
		return err
	}

	return nil
}

func RunLocalMinimum() {
	var N int          // Length of range
	var Out string     // Output file
	var noPrint bool   // Do not print intermediate steps
	var noVisuals bool // Do not visualize the graph
	var noSave bool    // Disable save prompt

	var seed = time.Now().UnixNano()
	var intSeed int

	// Extract args
	for i, v := range os.Args {
		switch v {
		case "-N", "--N":
			SScanInt(os.Args[i+1], &N, "N")
		case "-S", "--S":
			SScanInt(os.Args[i+1], &intSeed, "seed")
		case "-O", "--O":
			fmt.Sscanf(os.Args[i+1], "%s", &Out)
		case "--noprint":
			noPrint = true
		case "--novisuals":
			noVisuals = true
		case "--nosave":
			noSave = true
		}
	}

	fmt.Println("---- Local minimum finding program ----")

	if N == 0 {
		fmt.Print("Input the length of the range N: ")
		ScanInt(&N, "N")
	}

	if N < 3 {
		fmt.Print("Range cannot be shorter than 3")
		os.Exit(1)
	}

	if intSeed > 0 {
		seed = int64(intSeed)
	}

	fmt.Printf("---- Creating range of size %d with layered sine waves ----\n", N)
	fmt.Printf("Seed: %v\n", seed)

	r := rand.New(rand.NewSource(seed))
	floatRange := make([]float64, N)

	// Sine coefficients
	cof := []float64{
		r.Float64(),
		1 + r.Float64()*3,
		3 + r.Float64()*3,
		5 + r.Float64()*5,
		10 + r.Float64()*10,
		1.0 + r.Float64()*5.0,
		1.0 + r.Float64()*1.5,
		1.0 + r.Float64()*0.5,
		1.0 + r.Float64()*0.25}

	// Layered sine function generator
	generator := func(i int) float64 {
		x := math.Pi * (float64(i)/float64(N) + cof[0])
		s1 := math.Sin(x * cof[1])
		s2 := math.Sin(x * cof[2])
		s3 := math.Sin(x * cof[3])
		s4 := math.Sin(x * cof[4])
		return cof[5]*s1 + cof[6]*s2 + cof[7]*s3 + cof[8]*s4
	}

	for i := range floatRange {
		floatRange[i] = generator(i)
	}

	// Fix edge conditions to ensure a local minimum always exists
	floatRange[0] = floatRange[1] + 0.5
	floatRange[len(floatRange)-1] = floatRange[len(floatRange)-2] + 0.5

	// Find minimum
	localMinimum := FindLocalMinimum(floatRange, noPrint)
	fmt.Printf("Local minimum is %.2f at i = %v\n", floatRange[localMinimum], localMinimum)

	if !noSave && len(Out) < 1 {
		saveRange := "n"
		fmt.Println("Save range to file? (y/n)")
		fmt.Scanf("%s\n", &saveRange)
		if saveRange == "Y" {
			fmt.Print("Filename: ")
			fmt.Scanf("%s\n", &Out)
		}
	}

	if !noSave && len(Out) > 0 {
		if err := SaveRangeJson(floatRange, Out); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Saved to %s\n", Out)
		}
	}

	if noVisuals {
		fmt.Println("--novisuals specified")
	} else {
		visualize := "n"
		fmt.Println("Display the graph in a browser? (y/n)")
		fmt.Scanf("%s\n", &visualize)
		if visualize == "y" {
			PlotFunction(floatRange, localMinimum)
		}
	}
}
