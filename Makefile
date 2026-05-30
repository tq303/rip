.PHONY: install build clean release

install:
	go install .

build:
	go build -o rip .

clean:
	rm -f rip

release:
	./scripts/release.sh
