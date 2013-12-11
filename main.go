/*
Chunk node server code
SilentFl
2013
*/
package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strconv"
)

type ReadQuery struct {
		ChunkHandle int64
		StartIndex, EndIndex int
}

type Result struct {
		Data []byte
}

var (
	storagepath = flag.String("D", "/tmp", "path to chunkfile directory")
)

func (p *ReadQuery) Get(args ReadQuery, reply *Result) error {
	log.Println("get ", args.ChunkHandle)
	pathToChunk := filepath.Join(*storagepath, "chunk00"+strconv.FormatInt(args.ChunkHandle, 10))+".txt"
	log.Println("path to chunk: ", pathToChunk)
	chunk, err := os.Open(pathToChunk)
	if err != nil {
		return errors.New("Chunk not found")
	} else {
		r := make([]byte, args.EndIndex-args.StartIndex)
		chunk.ReadAt(r, int64(args.StartIndex))
		reply.Data = r
//		log.Println(r)
		return nil
	}
}

func main() {
	get := new(ReadQuery)
	rpc.Register(get)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":4001")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
