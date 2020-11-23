test:
	go test ./... -race -v
proto:
	protoc -I=./business/proto --go_out=plugins=grpc:../../.. ./business/proto/networking/networking.proto