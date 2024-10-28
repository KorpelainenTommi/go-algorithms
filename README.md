# go algorithms
Implementations for some challenge problems in go

- [X] Graph k-coloring 
- [X] Local minimum find
- [ ] Largest contiguous submatrix

## Prerequisites
- Install [Go](https://go.dev/), you should be able to run `go version`
- Run the program with `go run .`
## Args

By default, the program will ask for arguments to the problem parameters. These can also be provided as command line arguments. Here is a quick list of accepted arguments
- `-P` or `--P` Choose algorithm to demo
- `--noprint` Disable printing algorithm solutions / steps (useful if the solution is very long, e.g. Graph coloring)
- `--nosave` Disable save prompt
- `--novisuals` Disable prompt for visuals
- `-O` or `--O` Save output to file, e.g. `-O output.json`

### Graph coloring
- `-N` or `--N` number of nodes
- `-D` or `--D` desired average edge degree
- `-S` or `--S` seed for random generation

### Local minimum
- `-N` or `--N` length of the range
- `-S` or `--S` seed for random generation



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