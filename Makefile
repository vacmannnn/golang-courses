build:
	@go build -C cmd/xkcd -o "../../xkcd-server"

bench:
	@go test -bench=. -v ./cmd/xkcd/

clean:
	@rm xkcd

test:
	@go test ./... -v -race -cover -coverprofile coverage/coverage.out ## TODO: -race
	@go tool cover -html coverage/coverage.out -o coverage/coverage.html

lint:
	@golangci-lint run ./...

sec:
	@trivy fs xkcd-server
	@govulncheck ./...

e2e:
	./e2e.sh