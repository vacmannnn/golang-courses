build:
	@go build -C cmd/xkcd -o "../../xkcd-server"

bench:
	@go test -bench=. -v ./cmd/xkcd/

clean:
	@rm xkcd

test:
	@go test ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out ## TODO: -race
	@go tool cover -html coverage/coverage.out -o coverage/coverage.html
