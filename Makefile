start:
	@if [ ! -d src ]; then \
		git clone -b v2 https://github.com/zhengkai/coral.git src; \
	fi
	go run *.go

tag:
	./tag.sh

lint:
	cd src && golangci-lint run ./...
