go:
	go run main.go

gotest:
	go test ./... -cover

build:
	docker build -t scale-go .

run:
	docker run --rm scale-go

tests:
	docker build -t scale-test . -f Dockerfile.test
	docker run --rm scale-test

chglog:
	docker run --rm -v $(PWD):/workdir quay.io/git-chglog/git-chglog -o CHANGELOG.md