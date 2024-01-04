# This is a justfile. See https://github.com/casey/just

test:
	go test ./...

test-force:
	go test -count=1 ./...