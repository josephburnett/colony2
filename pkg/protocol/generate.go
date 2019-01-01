package protocol

//go:generate protoc --proto_path=$GOPATH/src:. --gogo_out=plugins=grpc:. ./colony.proto
