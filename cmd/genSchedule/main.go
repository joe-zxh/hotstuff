package main

import "fmt"

var clusterSize int = 160

func genSchedule() {
	fmt.Printf("[")
	for i := 1; i < clusterSize; i++ {
		fmt.Printf("%d, ", i)
	}
	fmt.Printf("%d]\n", clusterSize)
}

func genReplicas() {
	for i := 1; i <= clusterSize; i++ {
		fmt.Printf("[[replicas]]\n")
		fmt.Printf("id = %d\n", i)
		fmt.Printf("peer-address = \"127.0.0.1:%d\"\n", 13370+i)
		fmt.Printf("client-address = \"127.0.0.1:%d\"\n", 23370+i)
		fmt.Printf("pubkey = \"keys/r%d.key.pub\"\n", i)
		fmt.Printf("cert = \"keys/r%d.crt\"\n\n", i)
	}
}

func main() {
	genReplicas()
}
