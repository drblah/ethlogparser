default: bin bin/windows bin/macos
	GOARCH=amd64 GOOS=darwin go build -o bin/macos/ethlogparser
	GOARCH=amd64 GOOS=windows go build -o bin/windows/ethlogparser.exe

bin:
	mkdir bin

bin/windows:
	mkdir bin/windows

bin/macos:
	mkdir bin/macos

