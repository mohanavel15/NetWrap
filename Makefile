.PHONY: client
client:
	go build cmd/client/main.go
	mv main build/client

.PHONY: server
server:
	go build cmd/server/main.go
	mv main build/server
.PHONY: clean
clean:
	rm build/*
