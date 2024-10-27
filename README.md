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