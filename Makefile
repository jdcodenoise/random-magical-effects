
run: build
	./bin/magical-effects

build: bin
	go build -o ./bin/magical-effects

bin:
	mkdir ./bin
