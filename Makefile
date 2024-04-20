build:
	@go build -C cmd/xkcd -o "../../xkcd"

clean:
	rm xkcd
