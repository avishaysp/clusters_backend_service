# compiling IDL into Go and gRPC to to generate `RegService.pb.go` and `RegService_grpc.pb.go`
cd /workspaces/RLAD/services/reg-service/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MRegService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/reg-service \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=MRegService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/RegService.proto \
       RegService.proto