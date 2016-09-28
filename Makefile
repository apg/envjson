envjson:
	go build ./...

test: envjson
	prove --exec '/bin/sh' t/

clean:
	rm -f ./envjson

.PHONY: test clean
