build:
	@go build -C cmd/xkcd -o "../../xkcd-server"
	@go build -C cmd/web -o "../../web-server"

xkcd:
	@go build -C cmd/xkcd -o "../../xkcd-server"

web:
	@go build -C cmd/web -o "../../web-server"

run:
	./xkcd-server & ./web-server

bench:
	@go test -bench=. -v ./cmd/xkcd/

clean:
	@rm xkcd-server
	@rm web-server

test:
	@go test ./... -v -race -cover -coverprofile coverage/coverage.out
	@go tool cover -html coverage/coverage.out -o coverage/coverage.html

lint:
	@golangci-lint run ./...

sec:
	@trivy fs xkcd-server
	@govulncheck ./...

e2e:
	./e2e.sh