build:
	go build -ldflags="-w -s" twmd.go

all: build

clean:
	rm -f twmd
