start:
	go test -race -v -test.bench="test" -test.benchmem -coverpkg github.com/zhengkai/coral -covermode=atomic -coverprofile=coverage.out

html:
	go tool cover -html=coverage.out
