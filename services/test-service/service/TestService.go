package TestService

import (
	"context"
	"fmt"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	services "github.com/TAULargeScaleWorkshop/RLAD/services/common/common-service"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common" // from test-service/common import *
	TestServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/servant"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils" // from utils import *

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type testServiceImplementation struct {
	UnimplementedTestServiceServer
}

type mockWaitAndRandServer struct {
	grpc.ServerStream
	sendFunc func(msg *wrapperspb.Int32Value) error
}

func (m *mockWaitAndRandServer) Send(msg *wrapperspb.Int32Value) error {
	return m.sendFunc(msg)
}

var serviceInstance *testServiceImplementation

func messageHandler(method string, parameters []byte) (response proto.Message, err error) {
	Logger.Printf("messageHandler(): entered with %s and %v", method, parameters)
	switch method {
	case "ExtractLinksFromURL":
		p := &ExtractLinksFromURLParameters{}
		err := proto.Unmarshal(parameters, p)
		if err != nil {
			return nil, err
		}
		res, err := serviceInstance.ExtractLinksFromURL(context.Background(), p)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "Get":
		p := &wrapperspb.StringValue{}
		err := proto.Unmarshal(parameters, p)
		if err != nil {
			return nil, err
		}
		res, err := serviceInstance.Get(context.Background(), p)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "Store":
		p := &StoreKeyValue{}
		err := proto.Unmarshal(parameters, p)
		if err != nil {
			return nil, err
		}
		res, err := serviceInstance.Store(context.Background(), p)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "HelloWorld":
		p := emptypb.Empty{}
		res, err := serviceInstance.HelloWorld(context.Background(), &p)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "HelloToUser":
		p := &wrapperspb.StringValue{}
		err := proto.Unmarshal(parameters, p)
		if err != nil {
			return nil, err
		}
		res, err := serviceInstance.HelloToUser(context.Background(), p)
		if err != nil {
			return nil, err
		}
		return res, nil
	case "WaitAndRand":
		return nil, fmt.Errorf("WaitAndRand is unsupported")
	case "IsAlive":
		p := emptypb.Empty{}
		res, err := serviceInstance.IsAlive(context.Background(), &p)
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		return nil, fmt.Errorf("MQ message called unknown method: %v", method)
	}
}

func Start(configData []byte) error {
	// get base config
	var baseConfig config.BaseConfig
	err := yaml.Unmarshal(configData, &baseConfig) // parses YAML
	if err != nil {
		Logger.Fatalf("error unmarshaling BaseConfig data: %v", err)
	}

	// get TestService config
	var config config.TestConfig
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		Logger.Fatalf("error unmarshaling TestConfig data: %v", err)
	}
	config.BaseConfig = baseConfig

	TestServiceServant.InitServant(config.RegistryAddresses)

	// start service
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterTestServiceServer(s, serviceInstance)
	}
	services.Start(config.Type, 0, config.RegistryAddresses, bindgRPCToService, messageHandler) // randomly pick a port

	return nil
}

func (obj *testServiceImplementation) HelloWorld(ctxt context.Context, _ *emptypb.Empty) (*wrapperspb.StringValue, error) {
	Logger.Printf("HelloWorld called")
	return wrapperspb.String(TestServiceServant.HelloWorld()), nil
}

func (obj *testServiceImplementation) HelloToUser(_ context.Context, userName *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	Logger.Printf("HelloToUser called")
	return wrapperspb.String(TestServiceServant.HelloToUser(userName.Value)), nil
}

func (obj *testServiceImplementation) Store(ctx context.Context, req *StoreKeyValue) (*emptypb.Empty, error) {
	Logger.Printf("Store called with key: %s, value: %s", req.Key, req.Value)
	err := TestServiceServant.Store(req.Key, req.Value)
	return &emptypb.Empty{}, err
}

func (obj *testServiceImplementation) Get(ctx context.Context, key *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	Logger.Printf("Get called with key: %s", key.Value)
	value, err := TestServiceServant.Get(key.Value)
	return wrapperspb.String(value), err
}

func (obj *testServiceImplementation) WaitAndRand(seconds *wrapperspb.Int32Value, streamRet TestService_WaitAndRandServer) error {
	Logger.Printf("WaitAndRand called")
	streamClient := func(x int32) error {
		return streamRet.Send(wrapperspb.Int32(x))
	}
	return TestServiceServant.WaitAndRand(seconds.Value, streamClient)
}

func (obj *testServiceImplementation) IsAlive(ctxt context.Context, _ *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	Logger.Printf("IsAlive called")
	return wrapperspb.Bool(TestServiceServant.IsAlive()), nil
}

func (obj *testServiceImplementation) ExtractLinksFromURL(ctx context.Context, url *ExtractLinksFromURLParameters) (*ExtractLinksFromURLReturnedValue, error) {
	Logger.Printf("ExtractLinksFromURL called with url: %s", url.Url)
	value, err := TestServiceServant.ExtractLinksFromURL(url.Url, url.Depth)
	return &ExtractLinksFromURLReturnedValue{Links: value}, err
}
