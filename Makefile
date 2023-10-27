build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o bin/alertmanager2