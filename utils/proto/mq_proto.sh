# compiling IDL into Go and gRPC to to generate `RegService.pb.go` and `RegService_grpc.pb.go`
cd /workspaces/RLAD/services/common
protoc -I=. \
       --go_out=. \
       --go_opt=paths=source_relative \
       --go_opt=MCallMessage.proto=github.com/TAULargeScaleWorkshop/RLAD/services/common \
       CallMessage.proto