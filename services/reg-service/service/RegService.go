package RegService

import (
	"context"
	"fmt"
	"net"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
	RegServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/servant"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type regServiceImplementation struct {
	UnimplementedRegServiceServer
}

// Note: Code Duplication to prevent importing ServiceBase
func startgRPC(listenPort int) (listeningAddress string, grpcServer *grpc.Server, startListening func(), err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		// we treat this as "warning"
		Logger.Printf("startgRPC(): failed to listen: %v", err)
		return "", nil, nil, err
	}
	listeningAddress = lis.Addr().String()
	grpcServer = grpc.NewServer()
	startListening = func() {
		if err := grpcServer.Serve(lis); err != nil {
			Logger.Fatalf("failed to serve: %v", err)
		}
	}
	return
}

// Note: copy of ServiceBase::Start() without the regAddresses parameter (not needed)
func startRegService(name string, grpcListenPort int, bindgRPCToService func(s grpc.ServiceRegistrar)) (err error) {
	// start the service
	listeningAddress, grpcServer, startListening, err := startgRPC(grpcListenPort)
	if err != nil {
		return err
	}
	Logger.Printf("startRegService(): calling bindgRPCToService")
	bindgRPCToService(grpcServer)

	// init only when a new RegService is starting
	RegServiceServant.InitServant(name)

	// IsAlive check performed only by the "root" node
	if RegServiceServant.IsFirst() {
		Logger.Printf("RegService starts an IsAlive thread")
		go RegServiceServant.IsAliveCheck()
	}

	Logger.Printf("RegService starts listening on %s\n", listeningAddress)
	// blocks
	startListening()
	return nil
}

func Start(configData []byte) error {
	// get config
	var config config.RegConfig
	err := yaml.Unmarshal(configData, &config)
	if err != nil {
		Logger.Fatalf("error unmarshaling data: %v", err)
	}

	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterRegServiceServer(s, &regServiceImplementation{})
	}
	for port := config.ListenPort; port < config.ListenPort+10; port++ {
		Logger.Printf("Start(): trying port %d", port)
		err = startRegService(config.Name, port, bindgRPCToService)
		// will reach here only if failed to connect
		if err != nil {
			// we treat this as "warning"
			Logger.Printf("Start(): startRegService failed %v for port %d", err, port)
		}
	}
	return fmt.Errorf("failed to start RegService")
}

func (obj *regServiceImplementation) Register(_ context.Context, params *UpdateRegistryParameters) (_ *emptypb.Empty, err error) {
	Logger.Printf("Register called with service: %s, address: %s", params.ServiceName, params.NodeAddr)
	RegServiceServant.Register(params.ServiceName, params.NodeAddr)
	return &emptypb.Empty{}, nil
}

func (obj *regServiceImplementation) Unregister(_ context.Context, params *UpdateRegistryParameters) (_ *emptypb.Empty, err error) {
	Logger.Printf("Unregister called with service: %s, address: %s", params.ServiceName, params.NodeAddr)
	RegServiceServant.Unregister(params.ServiceName, params.NodeAddr)
	return &emptypb.Empty{}, nil
}

func (obj *regServiceImplementation) Discover(_ context.Context, service_name *wrapperspb.StringValue) (*DiscoveredServers, error) {
	Logger.Printf("Discover called with service: %s", service_name.Value)
	value, err := RegServiceServant.Discover(service_name.Value)
	return &DiscoveredServers{Addresses: value}, err
}
