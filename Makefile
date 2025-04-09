.PHONY: build dev test cov cov-html clean

build:
	go build .

dev: build
	./jaqen

test:
	go test ./...

cov:
	go test -coverprofile=/tmp/coverage.out ./...

cov-html: cov
	go tool cover -html=/tmp/coverage.out

clean:
	rm -f jaqen
	rm -f /tmp/coverage.out
