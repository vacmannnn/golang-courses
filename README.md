# xkcd searcher

A tool to search comics from [xkcd](https://xkcd.com/).

https://github.com/vacmannnn/xkcd-searcher/assets/111463436/c4cdb73d-8aba-4bd0-aa8d-024b661cdc93

## Usage

1. Download project:
```bash
git clone git@github.com:vacmannnn/xkcd-searcher.git && cd xkcd-searcher
```

2. Build and run it:
```bash
make
make run
```

3. Go to localhost:3000 to use search (default login and password is `user;user`)

If you want to run separate scenarios you may use:
- `make xkcd` to build `xkcd-server`
- `make web` to build `web-server`
- `make test` to run all tests
- `make lint/sec` to run utility checks on project (to run `make sec` you need to build `xkcd-server` first)

## web-server
Web-server connects with `xkcd-server` to use it as comics finder and authorization service.

Default source url is `localhost:3000/login`

## xkcd-server

By default, it listens `localhost:8080`. To use `xkcd-server` you may use rest API:

Getting JWT:
```bash
curl -d '{"username":"admin", "password":"admin"}' "http://localhost:8080/login"
```
Update and search requests:
```bash
curl -X POST --header "Authorization:TOKENFROMPREVIOUSSTEP" "http://localhost:8080/update"
```
```bash
curl -X GET --header "Authorization:TOKENFROMPREVIOUSSTEP" "http://localhost:8080/pics?search='apple,doctor'"
```
