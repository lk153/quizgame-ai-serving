GOPATH              := $(or $(GOPATH), $(HOME)/go)
GOLINT              := golangci-lint run --timeout 10m
MOCKERY             := $(GOPATH)/bin/mockery
GO_TEST_PARALLEL    := go test -parallel 4 -count=1 -timeout 30s
GOSTATIC            := go build -ldflags="-w -s"
GOOGLE_WIRE 		:= $(GOPATH)/bin/wire

$(MOCKERY):
	GOPATH=$(GOPATH) go install -mod=mod github.com/vektra/mockery/v2@latest
$(GOOGLE_WIRE):
	GOPATH=$(GOPATH) go install github.com/google/wire/cmd/wire@latest
start:
	./out/main
clean:
	rm -rf ./out/main cpu.pprof mem.pprof
build: $(GOOGLE_WIRE) clean
	go mod tidy && go mod vendor && $(GOOGLE_WIRE) ./cmd/http && $(GOSTATIC) -o out/main ./cmd/http
lint:
	$(GOLINT) -v ./...
test:
	$(GO_TEST_PARALLEL) ./... -v -coverprofile=cover.out && go tool cover -html=cover.out
test-race:
	$(GO_TEST_PARALLEL) ./... -v -failfast -count 1 -race -coverprofile=cover.out && go tool cover -html=cover.out
test-mono: #tc: testcase name | ft: folder test
	go test -run $(tc) $(ft) -v -failfast -count 1 -timeout 30s
mock: $(MOCKERY) #if: Interface will mock | dir: folder of interface | sn: name of mock struct
	$(MOCKERY) --name=$(if) --dir=$(dir) --structname=$(sn) --output=$(dir)/mocks
generate:
	go generate ./...