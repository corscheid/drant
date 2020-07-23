drant: format all

format:
	gofmt -w drant.go

all:
	go build -o bin/drant ./...

install: drant
	sudo cp bin/drant /usr/local/bin/

uninstall:
	sudo rm /usr/local/bin/drant

clean:
	rm bin/drant
