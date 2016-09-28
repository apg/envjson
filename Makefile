envjson: main.go env.go
	go build ./...

gotest: env.go env_test.go
	go test -v ./...

funtest: envjson
	prove --exec '/bin/sh' t/

test: gotest funtest

clean:
	rm -f ./envjson

.PHONY: test gotest funtest clean
