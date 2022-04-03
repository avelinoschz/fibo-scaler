root:
	curl localhost:8080/

current:
	curl localhost:8080/current

next:
	curl localhost:8080/next

prev:
	curl localhost:8080/previous

go:
	go run main.go

gobuild:
	go build -o bin/main

gotest:
	go test ./... -cover

build:
	docker build -t fibo-scaler .

run:
	docker run --rm fibo-scaler

tests:
	docker build -t fibo-scaler-tests . -f Dockerfile.test
	docker run --rm fibo-scaler-tests

chglog:
	docker run --rm -v $(PWD):/workdir quay.io/git-chglog/git-chglog -o CHANGELOG.md