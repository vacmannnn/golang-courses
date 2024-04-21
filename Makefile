build:
	@go build -C cmd/xkcd -o "../../xkcd"

test:
	@go test -bench=. -v ./cmd/xkcd/

clean:
	@rm xkcd
