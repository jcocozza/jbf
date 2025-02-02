bin=jbf

build:
	go build -o $(bin) cmd/main.go

dev:
	go build -o $(bin) cmd/main.go
	mv $(bin) test
