NAME = crawler

all: deps $(NAME)

${NAME}:
	go build -o $(NAME) 

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

clean:
	rm -rf $(NAME) vendor/* dist coverage.out
	go clean -i ./...