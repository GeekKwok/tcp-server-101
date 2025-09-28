all: server client

server: cmd/server/main.go
	go build github.com/geekkwok/tcp-server-101/cmd/server
client: cmd/client/main.go
	go build github.com/geekkwok/tcp-server-101/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client