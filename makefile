build:
	go build

run:
	go build
	./voting

test1:
	go build
	./voting -mode test -i 1

test2:
	go build
	./voting -mode test -i 2

test:
	go build
	./voting -mode test