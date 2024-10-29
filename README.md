# go algorithms
Implementations for some challenge problems in go

- [X] Graph k-coloring
- [X] Local minimum find
- [X] Largest contiguous submatrix (although it is quite slow)

## Prerequisites
- Install [Go](https://go.dev/), you should be able to run `go version`
- Run the program with `go run .`

## Args
By default, the program will ask for arguments to the problem parameters. These can also be provided as command line arguments. Here is a quick list of accepted arguments
- `-P` Choose algorithm to demo
- `-O` Save output to file, e.g. `-O output.json`
- `-S` Seed for random generation
- `--noprint` Disable printing algorithm solutions / steps (useful if the solution is very long, e.g. Graph coloring)
- `--nosave` Disable save prompt
- `--novisuals` Disable prompt for visuals

### Graph coloring
- `-N` Number of nodes
- `-D` Desired average edge degree. Floating point number.

### Local minimum
- `-N` Length of the range

### Submatrix
- `-N` Image/matrix dimensions NxN
- `-B` Image blockiness, higher is blockier. Positive integer
- `-I` Path to an input image


## Example usage
```
go run . --noprint
==== Algorithm demo ====
Choose the program to execute:
1. Graph coloring
2. Local minimum
3. Largest contiguous submatrix
Program: 1
---- Graph coloring program ----
Input the number of nodes N: 10000
Input average node degree D: 4
---- Creating graph of size 10000 with average node degree < 4.00 ----
Seed: 1730144250400212500
Connected: true
Reachable: 10000
Colored: true
Colors used: 5
Successfully colored the graph
--noprint specified, skip graph printing
Save graph to file? (y/n)
y
Filename: graph.txt
Saved to graph.txt
Display the graph in a browser? (y/n)
n
```

------

```
go run . --nosave --novisuals
==== Algorithm demo ====
Choose the program to execute:
1. Graph coloring
2. Local minimum
3. Largest contiguous submatrix
Program: 2
---- Local minimum finding program ----
Input the length of the range N: 800
---- Creating range of size 800 with layered sine waves ----
Seed: 1730144983795741600
start: 0 end: 800 i: 400
start: 0 end: 400 i: 200
start: 200 end: 400 i: 300
start: 200 end: 300 i: 250
start: 250 end: 300 i: 275
start: 275 end: 300 i: 287
start: 287 end: 300 i: 293
start: 287 end: 293 i: 290
start: 287 end: 290 i: 288
start: 288 end: 290 i: 289
Local minimum is -3.38 at i = 289
--novisuals specified
```

------

```
go run . -P 3 -I "./bw.png" --noprint --novisuals
==== Algorithm demo ====
---- Largest contiguous submatrix ----
Attempting to read input file: ./bw.png
Successfully read png file
Record is the area rectangle (0,67) -> (73,128)
--novisuals specified
```

------

![Graph color image](https://github.com/KorpelainenTommi/go-algorithms/blob/main/graph-color-example.png)

------

![Local minimum image](https://github.com/KorpelainenTommi/go-algorithms/blob/main/local-minimum-example.png)

------

![Submatrix example 1](https://github.com/KorpelainenTommi/go-algorithms/blob/main/bw-output.png)

------

![Submatrix example 2](https://github.com/KorpelainenTommi/go-algorithms/blob/main/256-example.png)