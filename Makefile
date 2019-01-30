# Go parameters
GOBUILD= go  build
GOCLEAN= go  clean
GOTEST= go  test
BINARY_NAME=EndoRestAPI
all: test build
build: 
	$(GOBUILD) -ldflags "-X main.VersionString=`git rev-parse HEAD`" -o $(BINARY_NAME) -v ./...
test: 
	$(GOTEST) -v  ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
docker-build:
	docker run --rm -it -v $(PWD):/go -w /go golang:latest  go test && go build -ldflags "-X main.VersionString=`git rev-parse HEAD`"  -o "$(BINARY_NAME)" && ./EndoRestAPI
