test:
	go test ./... -race -v
proto:
	protoc -I=./foundation/proto --go_out=plugins=grpc:../../.. ./foundation/proto/networking/networking.proto ./foundation/proto/messaging/messaging.proto