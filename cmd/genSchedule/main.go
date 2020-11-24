package main

import "fmt"

func main() {
	clusterSize := 52
	fmt.Printf("[")
	for i := 1; i < clusterSize; i++ {
		fmt.Printf("%d, ", i)
	}
	fmt.Printf("%d]\n", clusterSize)
}
