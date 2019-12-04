build:
	 env GOOS=darwin go build -o bin/darwin/lambda-local
	 env GOOS=linux  go build -o bin/linux/lambda-local
install:
	go install

