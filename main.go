package main

import (
	"../shared"
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

var (
	storagepath = flag.String("/tmp", "D", "path to chunkfile directory")
)

type localReadQuery shared.ReadQuery

func (p *localReadQuery) Get(args *shared.ReadQuery, reply []byte) error {
	log.Println("get ", args.ChunkHandle)
	chunk, err := os.Open(filepath.Join(*storagepath, strconv.FormatInt(args.ChunkHandle, 10)))
	if err != nil {
		return errors.New("Chunk not found")
	} else {
		reply = make([]byte, args.EndIndex-args.StartIndex)
		chunk.ReadAt(reply, int64(args.StartIndex))
		return nil
	}
}

func main() {
	get := new(localReadQuery)
	rpc.Register(get)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":4001")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
