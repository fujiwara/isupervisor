.PHONY: clean test

isupervisor: go.* *.go cmd/isupervisor/*.go
	go build -o $@ cmd/isupervisor/main.go

clean:
	rm -rf isupervisor dist/

test:
	go test -v ./...

install:
	go install github.com/fujiwara/isupervisor/cmd/isupervisor

dist:
	goreleaser build --snapshot --rm-dist
