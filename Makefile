prog=urlshortener

all:
	go build -o $(prog) cmd/api/*

test:
	go test -count=1 -v ./...

clean:
	rm -rf $(prog)