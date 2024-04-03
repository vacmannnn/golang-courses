build:
	@go build -o "xkcd" cmd/xkcd/main.go

clean:
	rm xkcd