lint:
	go mod tidy
	go fmt ./...
	go vet ./...
	golangci-lint run ./...

pb:
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative example/pb/*.proto
	frgen pbfinger --in=example/pb/hello.proto
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative test/full/pb/*.proto
	frgen pbfinger --in=test/full/pb/hello.proto --nats=true
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative test/dual/pb/*.proto
	frgen pbfinger --in=test/dual/pb/hello.proto --nats=true
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative test/get/pb/*.proto
	frgen pbfinger --in=test/get/pb/hello.proto --nats=true
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative test/set/pb/*.proto
	frgen pbfinger --in=test/set/pb/hello.proto --nats=true
	protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative test/touch/pb/*.proto
	frgen pbfinger --in=test/touch/pb/hello.proto --nats=true
