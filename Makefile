clean:
	rm -vf bin/*
build: clean
	go build -o bin/dumpster dumpster.go
