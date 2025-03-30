build: 
	go build .

dev:
	make build && ./jaqen

cov:
	go test  -coverprofile=/tmp/coverage.out ./...

cov-html:
	make cov
	go tool cover -html=/tmp/coverage.out
