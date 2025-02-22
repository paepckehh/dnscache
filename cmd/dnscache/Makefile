all:
	GOGC="off" go test -bench . 

check: 
	gofmt -w -s .
	staticcheck
