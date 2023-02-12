# Protobuf generated go files
PROTO_GO_FILES = $(shell find . -path -prune -o -type f -name '*.pb.go' -print | grep -v vendor)
#PROTO_GO_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))
DEST=${PWD}
BIN=trustedengine

.PHONY: all bin proto build deps clean sim run host

all: proto bin

bin:
	ego-go build -o=${BIN} ./cmd/trustedengine
	ego sign ${BIN}

build: $(PROTO_GO_FILES)
	@buf build

sim:
	OE_SIMULATION=1 ego run ${BIN}

run:
	ego run ${BIN}

host:
	go build -o=host_${BIN} ./cmd/trustedengine

proto: build
	@buf generate

deps:
	@go get github.com/gogo/protobuf/proto
	@go install github.com/gogo/protobuf/protoc-gen-gogo 
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	@go install github.com/bufbuild/buf/cmd/buf@v1.9.0
	
clean:
	@rm -f $(PROTO_GO_FILES) $(BIN)
