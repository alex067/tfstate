compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm go build -o bin/main-darwin-arm main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/main-darwin-arm64 main.go
	GOOS=darwin GOARCH=386 go build -o bin/main-darwin-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/main-windows-amd64 main.go