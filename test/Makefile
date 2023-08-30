start:
	go test -v -test.bench="test" -test.benchmem -coverpkg github.com/zhengkai/coral/v2 -covermode=atomic -coverprofile=coverage.out

start:
	go test -race -v -test.bench="test" -test.benchmem -coverpkg github.com/zhengkai/coral/v2 -covermode=atomic -coverprofile=coverage.out

html:
	go tool cover -html=coverage.out
