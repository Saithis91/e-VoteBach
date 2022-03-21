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

test3:
	go build
	./voting -mode test -i 3
	
test4:
	go build
	./voting -mode test -i 4

test:
	go build
	./voting -mode test