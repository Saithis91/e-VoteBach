build:
	go build

run:
	go build
	./voting

test1:
	go build
	./voting -mode test -i 1