install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

build_proto:
	@echo "Building Proto"
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    model/time.proto

build_server:
	@echo "Building Server"
	go build -o bin/server server/main.go

build_client:
	@echo "Building Client"
	go build -o bin/client client/main.go

run_server:
	@echo "Starting Server"
	./bin/server

run_client:
	@echo "Running Client"
	./bin/client