.PHONY: gen

gen:
	protoc --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:.  --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:.  --go_opt=paths=source_relative  --go-grpc_opt=paths=source_relative samerkv/samerkv.proto
