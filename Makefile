GOFILES=`find . -name "*.go" -type f`
PACKAGES=`go list ./...`
VETPACKAGES=`go list ./... | grep -v /examples/`

default:tidy fmt vet
	@go build -o main -v .

clean:
	@rm -f main
	@find . -name "[._]*.s[a-w][a-z]" | xargs -i rm -f {}

fmt:
	@gofmt -s -w ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

tidy:
	@go mod tidy

vet:
	@go vet ${VETPACKAGES}

test:
	@go test -v -count=1  ./...

docker:
	@docker build -t muxih4ck/Go-Web-Application-Template:latest .

ca:
	openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=xxxxx@qq.com"

help:
	@echo "make - compile the source code and tidy fmt vet"
	@echo "make fmt - format the source code"
	@echo "make fmt-check - check the source code"
	@echo "make tidy - go mod tidy"
	@echo "make vet - vet the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make docker - build a docker image"
	@echo "make ca - generate ca files"

.PHONY: default clean fmt fmt-check tidy vet test docker ca help
