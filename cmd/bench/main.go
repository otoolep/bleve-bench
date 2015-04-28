package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/otoolep/bleve-bench"
)

var batchSize = flag.Int("batchSize", 100, "batch size for indexing")
var nShards = flag.Int("shards", 1, "number of indexing shards")
var maxprocs = flag.Int("maxprocs", 1, "GOMAXPROCS")
var indexPath = flag.String("index", "indexes", "index storage path")
var docsPath = flag.String("docs", "docs", "path to docs file")
var csv = flag.Bool("csv", false, "summary CSV output")

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*maxprocs)

	// Remove any existing indexes.
	if err := os.RemoveAll(*indexPath); err != nil {
		fmt.Println("failed to remove %s.", *indexPath)
		os.Exit(1)
	}

	// Attempt to open the file.
	fmt.Printf("Opening docs file %s\n", *docsPath)
	f, err := os.Open(*docsPath)
	if err != nil {
		fmt.Printf("failed to open docs file: %s\n", err.Error())
		os.Exit(1)
	}

	// Read the lines into memory.
	docs := make([][]byte, 0, 100000)
	reader := bufio.NewReader(f)

	var l []byte
	l, err = reader.ReadBytes(byte('\n'))
	for err == nil {
		docs = append(docs, l)
		l, err = reader.ReadBytes(byte('\n'))
	}
	fmt.Printf("%d documents read for indexing.\n", len(docs))

	if len(docs)%(*nShards) != 0 {
		fmt.Println("Document count must be evenly divisible by shard count")
		os.Exit(1)
	}

	i := indexer.New(*indexPath, *nShards, *batchSize)
	if err := i.Open(); err != nil {
		fmt.Println("failed to open indexer:", err)
		os.Exit(1)
	}

	startTime := time.Now()
	if err := i.Index(docs); err != nil {
		fmt.Println("failed to index documents:", err)
		os.Exit(1)
	}
	duration := time.Now().Sub(startTime)

	count, err := i.Count()
	if err != nil {
		fmt.Println("failed to determine total document count")
		os.Exit(1)
	}
	rate := int(float64(count) / duration.Seconds())

	fmt.Printf("Commencing indexing. GOMAXPROCS: %d, batch size: %d, shards: %d.\n",
		runtime.GOMAXPROCS(-1), *batchSize, *nShards)

	fmt.Println("Indexing operation took", duration)
	fmt.Printf("%d documents indexed.\n", count)
	fmt.Printf("Indexing rate: %d docs/sec.\n", rate)

	if *csv {
		fmt.Printf("csv,%d,%d,%d,%d,%d,%d\n", len(docs), count, runtime.GOMAXPROCS(-1), *batchSize, *nShards, rate)
	}
}
