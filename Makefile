


build: 
	pigeon -o internal/parser/parser.go internal/parser/template.peg 
	go build -o bin/reason.exe

run: build 
	go run . 

