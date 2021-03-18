.PHONY: protos

protos:
	protoc -I protos/ --go-grpc_out=protos/currency --go_out=protos/currency protos/currency.proto