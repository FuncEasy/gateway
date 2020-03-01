.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOPROXY=https://goproxy.io GO111MODULE=on \
	go build -o ./build/bin/gateway -v -ldflags -s -a -installsuffix cgo ./main.go
	docker build --no-cache -t ziqiancheng/funceasy-gateway:latest .
	docker push ziqiancheng/funceasy-gateway:latest
clean:
	rm -rf ./build/bin