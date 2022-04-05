root:
	curl localhost:8080/

current:
	curl localhost:8080/current

next:
	curl localhost:8080/next

prev:
	curl localhost:8080/previous

err:
	curl localhost:8080/error

go:
	go run main.go

gobuild:
	go build -o bin/main

gotest:
	go test ./... -cover

build:
	docker build -t fibo-scaler .

run:
	docker run --rm -p 8080:8080 fibo-scaler

run-load:
	docker run -d -p 8080:8080 --restart=on-failure:3 --memory=512m --cpus=1  fibo-scaler

tests:
	docker build -t fibo-scaler-tests . -f Dockerfile.test
	docker run --rm fibo-scaler-tests

load:
	docker run --network=host --rm -it nakabonne/ali ali --rate 1000 --duration 1s http://localhost:8080/next

chglog:
	docker run --rm -v $(PWD):/workdir quay.io/git-chglog/git-chglog -o CHANGELOG.md