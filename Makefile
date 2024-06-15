build:
	env GOOS=linux GOARCH=amd64 go build -v -o reverse_me .

run:
	go run . --listen=127.0.0.1:8080 --target=http://127.0.0.1:8767
