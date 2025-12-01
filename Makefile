build-linux:
	GOOS=linux GOARCH=amd64 go build -tags=production -o bin/otari-linux-amd64 ./cmd/otari/
	GOOS=linux GOARCH=arm64 go build -tags=production -o bin/otari-linux-arm64 ./cmd/otari/