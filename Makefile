envjson: main.go env.go
	go build ./...

test: envjson
	prove --exec '/bin/sh' t/

clean:
	rm -f ./envjson

.PHONY: test clean
