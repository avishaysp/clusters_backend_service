# compiling IDL into Go and gRPC to to generate `TestService.pb.go` and `TestService_grpc.pb.go`
cd /workspaces/RLAD/services/test-service/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MTestService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/test-service \
       --go-grpc_out=. \
       --go-grpc_opt=paths=source_relative \
       --go-grpc_opt=MTestService.proto=github.com/TAULargeScaleWorkshop/RLAD/services/test-service/TestService.proto \
       TestService.proto