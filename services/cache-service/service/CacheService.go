package CacheService

import (
	"context"
	"fmt"

	"github.com/TAULargeScaleWorkshop/RLAD/config"
	. "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/common"
	CacheServiceServant "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/servant"
	services "github.com/TAULargeScaleWorkshop/RLAD/services/common/common-service"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v2"
)

type cacheServiceImplementation struct {
	UnimplementedCacheServiceServer
}

func Start(configData []byte) error {
	// get base config
	var baseConfig config.BaseConfig
	err := yaml.Unmarshal(configData, &baseConfig) // parses YAML
	if err != nil {
		Logger.Fatalf("error unmarshaling BaseConfig data: %v", err)
	}

	// get CacheService config
	var config config.CacheConfig
	err = yaml.Unmarshal(configData, &config) // parses YAML
	if err != nil {
		Logger.Fatalf("error unmarshaling CacheConfig data: %v", err)
	}
	config.BaseConfig = baseConfig

	// init only when a new CacheService is starting
	CacheServiceServant.InitServant(config.Name)

	// start service
	bindgRPCToService := func(s grpc.ServiceRegistrar) {
		RegisterCacheServiceServer(s, &cacheServiceImplementation{})
	}
	services.Start(config.Type, 0, config.RegistryAddresses, bindgRPCToService, nil) // randomly pick a port

	return nil
}

func (cs *cacheServiceImplementation) Set(_ context.Context, params *SetKeyValueReq) (_ *emptypb.Empty, err error) {
	Logger.Printf("Set called with key: %s, value: %s", params.Key, params.Value)
	CacheServiceServant.Set(params.Key, params.Value)
	return &emptypb.Empty{}, nil
}

func (cs *cacheServiceImplementation) Get(_ context.Context, k *GetKeyReq) (*GetValueReq, error) {
	Logger.Printf("Get called with key: %s", k.Key)
	value, err := CacheServiceServant.Get(k.Key)
	if err != nil {
		return nil, fmt.Errorf("key not found %s", k.Key)
	}
	return &GetValueReq{Value: value}, nil
}

func (cs *cacheServiceImplementation) Delete(_ context.Context, k *GetKeyReq) (_ *emptypb.Empty, err error) {
	Logger.Printf("Delete called with key: %s", k.Key)
	CacheServiceServant.Delete(k.Key)
	return &emptypb.Empty{}, nil
}

func (cs *cacheServiceImplementation) IsAlive(_ context.Context, _ *emptypb.Empty) (*wrapperspb.BoolValue, error) {
	Logger.Printf("IsAlive called ")
	IsAlive := CacheServiceServant.IsAlive()
	return wrapperspb.Bool(IsAlive), nil
}
