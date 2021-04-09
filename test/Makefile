start:
	go test -v -test.bench="test" -test.benchmem -coverpkg github.com/zhengkai/coral -covermode=count -coverprofile=coverage.out

html:
	go tool cover -html=coverage.out
