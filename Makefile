.PHONY: install build clean release

install:
	go build -o /usr/local/bin/rip .

build:
	go build -o rip .

clean:
	rm -f rip

release:
	./scripts/release.sh
