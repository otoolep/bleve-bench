bleve-bench
======

*Detailed background on this code can be found on [this blog post](http://www.philipotoole.com/increasing-bleve-performance-sharding/).*

bleve-bench is a program to test the impact of batch size and sharding on indexing performance of the [bleve library](https://github.com/blevesearch/bleve).

## Building and Running
*Building bleve-bench requires Go 1.3 or later. [gvm](https://github.com/moovweb/gvm) is a great tool for managing your version of Go.*

Download and run bleve-bench like so (tested on 64-bit Kubuntu 14.04):

    mkdir bleve-bench # Or any directory of your choice.
    cd bleve-bench/
    export GOPATH=$PWD
    go get -v github.com/otoolep/bleve-bench
    go install github.com/otoolep/bleve-bench/cmd/bench/.
    $GOPATH/bin/bench -h

Executing the last command will show the various options. An example run is shown below.

    $ $GOPATH/bin/bench -docs testdata.txt -maxprocs 8 -shards 50 -batchSize 100
    Opening docs file testdata.txt
    100000 documents read for indexing.
    Commencing indexing. GOMAXPROCS: 8, batch size: 100, shards: 50.
    Indexing operation took 3.479690221s
    100000 documents indexed.
    Indexing rate: 28738 docs/sec.
    
Each line in the test data file is read and indexed as a distinct document. Any previously indexed data is deleted before indexing begins.




