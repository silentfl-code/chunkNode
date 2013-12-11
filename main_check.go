package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type ReadQuery struct {
	ChunkHandle          int64
	StartIndex, EndIndex int
}

type Result struct {
	Data []byte
}

var reply Result

func main() {
	client, err := rpc.DialHTTP("tcp", "192.168.1.101:4001")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	fmt.Printf("Create rpc-connection: OK\n")
	// Synchronous call
	args := &ReadQuery{1, 0, 10}
	err = client.Call("ReadQuery.Get", args, &reply)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Printf("Get: %v\n", reply.Data)
}
