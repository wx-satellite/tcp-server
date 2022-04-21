all: server client

server: cmd/server/main.go
	go build ./cmd/server

client: cmd/client/main.go
	go build ./cmd/client

clean:
	rm -rf ./server
	rm -rf ./client

