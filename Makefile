GOBUILD=go build
GOTEST=go test

all:clean stop build_server
	./ecommerce 
build_server:
	$(GOBUILD) -v .
clean:
	rm -f ./ecommerce
stop:
	pkill ecommerce || true
test:
	cd helper && $(GOTEST) -v .