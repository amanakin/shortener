.PHONY: bin/shortener clean all test vendor

all: bin/shortener

bin/shortener:
	go build -mod=vendor -v -o bin/shortener cmd/shortener/main.go

protogen:
	protoc --proto_path=api/proto --go-grpc_out=internal/handler/grpc/api shortener.proto
	protoc --proto_path=api/proto --go_out=internal/handler/grpc/api shortener.proto

mockgen:
	mockgen -source=internal/repository/repository.go 	-destination=internal/mocks/mock_repository.go 	-package=mocks Repository
	mockgen -source=internal/service/service.go 		-destination=internal/mocks/mock_shortener.go 	-package=mocks Shortener
	mockgen -source=internal/service/shortener/shortener.go -destination=internal/mocks/mock_generator.go -package=mocks Generator

vendor:
	go mod vendor

test:
	go test -mod=vendor -v -race ./...

clean:
	rm -fv bin/shortener

