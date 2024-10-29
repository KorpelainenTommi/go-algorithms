package main

import (
	_ "encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"sync"
	"time"
)

type CounterMode uint8

const (
	Clean CounterMode = iota
	Dirty
	Black
	White
)

type CounterPool map[int]*Counter

type RecordMutex struct {
	mu sync.Mutex
	area int
	recordArea image.Rectangle
}

type Counter struct {
	id int
	finished bool
	in, out chan bool
	countArea image.Rectangle
}

func (c *Counter) BeginCount(recordMutex *RecordMutex) {

	var start, head, count int
	mode := Clean

	W := c.countArea.Dx()
	N := c.countArea.Dy()

	for head < N {
		
		value, ok := <- c.in
		if !ok {
			break
		}

		c.out <- true

		switch mode {
			case Clean:
				start = head
				if value {
					mode = Black
				} else {
					mode = White
				}
			case Black:
				if !value {
					mode = Dirty
				}
			case White:
				if value {
					mode = Dirty
				}
		}

		count++

		if count >= W {
			maxArea := W * (N - start)
			if maxArea <= recordMutex.area {
				<-c.in
				c.out <- false
				break
			}

			head++
			count = 0
			if mode == Dirty {
				mode = Clean
			} else {
				area := W * (head - start)
				if area > recordMutex.area {
					recordMutex.mu.Lock()
					if area > recordMutex.area {
						recordMutex.area = area
						recordMutex.recordArea.Min.X = c.countArea.Min.X
						recordMutex.recordArea.Max.X = c.countArea.Max.X
						recordMutex.recordArea.Min.Y = start
						recordMutex.recordArea.Max.Y = head
					}
					recordMutex.mu.Unlock()
				}
			}
		}
	}

	for _, ok := <-c.in; ok; _, ok = <-c.in {}
	close(c.out)
}

func FindLargest(N int, matrix []bool, reportProgress bool) image.Rectangle {
		
	recordMutex := new(RecordMutex)
	pools := make(map[int]CounterPool)
	counters := make([]Counter, 0)
	retired := make(map[int]bool)

	id := 0
	for w := 1; w <= N; w++ {
		for i := range N - w + 1 {
			v := Counter{
				id: id,
				in: make(chan bool),
				out: make(chan bool),
				countArea: image.Rectangle{
					Min: image.Point{i, 0},
					Max: image.Point{i+w, N},
				}}
			counters = append(counters, v)
			for vw := range w {
				if pools[i+vw] == nil {
					pools[i+vw] = make(CounterPool)
				}
				pools[i+vw][id] = &v
			}
			id++
		}
	}

	for c := range counters {
		go counters[c].BeginCount(recordMutex)
	}

	percent := 0
	for i, b := range matrix {
		v := pools[i % N]

		for _, counter := range v {
			counter.in <- b
			active, ok := <- counter.out

			if !active || !ok {
				for _, p := range pools {
					delete(p, counter.id)
				}
				close(counter.in)
				retired[counter.id] = true
			}
		}

		if reportProgress {
			p := 100 * i / (N*N)
			if p > percent {
				percent = p
				fmt.Println("Progress:",percent,"%")
			}
		}
	}

	for _, counter := range counters {
		if !retired[counter.id] {
			close(counter.in)
		}
	}

	// Block until all channels close
	for _, c := range counters {
		for _, ok := <-c.out; ok; _, ok = <-c.out {}
	}

	return recordMutex.recordArea
}

func ReadInputFile(path string) (int, []bool, error) {
	fmt.Printf("Attempting to read input file: %s \n", path)
	file, err := os.Open(path)
	if err != nil {
		return 0, nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return 0, nil, err
	}

	fmt.Println("Successfully read",format,"file")
	W, H := img.Bounds().Dx(), img.Bounds().Dy()

	if W != H {
		return 0, nil, errors.New(fmt.Sprintf("Image must be square, found: %v x %v", W, H))
	}

	L := W * H
	data := make([]bool, L)

	var thres uint32 = 0xFFFF / 4
	B := img.Bounds()
	x, y := B.Min.X, B.Min.Y

	for i := 0; i < L; i++ {
		r,g,b,_ := img.At(x, y).RGBA()
		data[i] = r <= thres && g <= thres && b <= thres
		x++
		if x >= B.Max.X {
			x = B.Min.X
			y++
		}
	}

	return W, data, nil

}

func RunSubmatrix() {
	var N int          // Width and height of the matrix NxN
	var B int          // Image blockiness
	var In string      // Input .png or .jpg file for the matrix
	var Out string     // Output file
	var noPrint bool   // Do not print intermediate steps
	var noVisuals bool // Do not visualize the graph
	var noSave bool    // Disable save prompt

	_ = noSave
	_ = noPrint

	var seed = time.Now().UnixNano()
	var intSeed int

	var matrix []bool // Black and white square matrix

	// Extract args
	for i, v := range os.Args {
		switch v {
		case "-N", "--N":
			SScanInt(os.Args[i+1], &N, "N")
		case "-B", "--B":
			SScanInt(os.Args[i+1], &B, "B")
		case "-S", "--S":
			SScanInt(os.Args[i+1], &intSeed, "seed")
		case "-O", "--O":
			fmt.Sscanf(os.Args[i+1], "%s", &Out)
		case "-I", "--I":
			fmt.Sscanf(os.Args[i+1], "%s", &In)
		case "--noprint":
			noPrint = true
		case "--novisuals":
			noVisuals = true
		case "--nosave":
			noSave = true
		}
	}

	fmt.Println("---- Largest contiguous submatrix ----")

	if len(In) == 0 && N == 0 {
		fmt.Print("Side length of the square matrix N: ")
		ScanInt(&N, "N")
	}

	if len(In) == 0 && B == 0 {
		fmt.Print("Image blockiness B (positive integer, higher means more blocky): ")
		ScanInt(&B, "B")
	}

	if intSeed > 0 {
		seed = int64(intSeed)
	}

	if len(In) > 0 {
		n, img, err := ReadInputFile(In)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		matrix = img
		N = n
	} else {
		fmt.Printf("---- Creating random blocky black and white image of size %dx%d ----\n", N, N)
		fmt.Printf("Seed: %v\n", seed)
		r := rand.New(rand.NewSource(seed))
		matrix = make([]bool, N*N)
		for i := range matrix {
			matrix[i] = r.Intn(1+B) > 0
		}
	}


	record := FindLargest(N, matrix, !noPrint)
	fmt.Printf("Record is the area rectangle (%v,%v) -> (%v,%v)\n",record.Min.X, record.Min.Y, record.Max.X, record.Max.Y)

	// if !noSave && len(Out) < 1 {
	// 	saveRange := "n"
	// 	fmt.Println("Save range to file? (y/n)")
	// 	fmt.Scanf("%s\n", &saveRange)
	// 	if saveRange == "Y" {
	// 		fmt.Print("Filename: ")
	// 		fmt.Scanf("%s\n", &Out)
	// 	}
	// }

	// if !noSave && len(Out) > 0 {
	// 	if err := SaveRangeJson(floatRange, Out); err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		fmt.Printf("Saved to %s\n", Out)
	// 	}
	// }

	if noVisuals {
		fmt.Println("--novisuals specified")
	} else {
		visualize := "n"
		fmt.Println("Display the graph in a browser? (y/n)")
		fmt.Scanf("%s\n", &visualize)
		if visualize == "y" {
			PlotImage(N, matrix, record)
		}
	}
}
