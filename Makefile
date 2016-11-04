build:
	go build .

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a .
	docker build --force-rm -t sh4rk/docker-socket-proxy .
	docker push sh4rk/docker-socket-proxy

.PHONY: build test end2end_test
