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

type ChunkNode bool

type ReadQuery struct {
		ChunkHandle int64
		StartIndex, EndIndex int
}

type Result struct {
		Data []byte
}

type WriteQuery struct {
		ChunkHandle int64
		Data []byte
}

var (
	storagepath = flag.String("D", "./tmp", "path to chunkfile directory")
	port = flag.Int("port", 4001, "Port for rpc-connection")
)

func init() {
	_, err := os.Stat(storagepath)
	if err != nil {
		if e, ok := err.(*os.PathError); ok && e.Error == os.ENOENT {
				err := os.Mkdir(".tmp", FileMode.ModeDir)
				if err != nil {
						log.Println("Error create dir for chunks!")
				}
		}
	}
}

func (p *ChunkNode) Get(args ReadQuery, reply *Result) error {
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

func  (p *ChunkNode) Write(args WriteQuery, success *bool) error {
	pathToChunk := filepath.Join(*storagepath, "chunk00"+strconv.FormatInt(args.ChunkHandle, 10))+".txt"
	file, err := os.Create(pathToChunk)
	defer file.Close()
	if err != nil {
			return errors.New("Error create chunk")
	} else {
			_, err := file.Write(args.Data)
			if err != nil {
					return errors.New("Error to write chunk")
			}
			return nil
	}
}

func main() {
	node := new(ChunkNode)
	rpc.Register(node)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":4001")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
