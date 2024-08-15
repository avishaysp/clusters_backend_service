# compiling IDL into Go and gRPC to to generate `RegService.pb.go` and `RegService_grpc.pb.go`
cd /workspaces/RLAD/services/cache-service/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MCacheService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/cache-service \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=MCacheService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/CacheService.proto \
       CacheService.proto