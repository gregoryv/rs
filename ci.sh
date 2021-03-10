#!/bin/bash -e

case $1 in

    bench)
	rm -rf *.out
	go test $2 \
	   -memprofile mem.out \
	   -run=Test_mem
	
	rm -rf *prof.png
	go tool pprof -png  -output memprof.png mem.out
	;;
    
    test)
	go test $2 -short -coverprofile /tmp/c.out ./...
	uncover /tmp/c.out
	;;
    
    *)
	$0 test
	;;
    
esac
