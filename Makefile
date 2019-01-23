drant: format all

format:
	gofmt -w drant.go

all:
	go build -o bin/drant drant.go

install: drant
	sudo cp bin/drant /usr/local/bin/

uninstall:
	sudo rm /usr/local/bin/drant

clean:
	rm bin/drant
