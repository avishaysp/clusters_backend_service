package CacheServiceClient

import (
	context "context"
	"fmt"

	service "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/common"
	client "github.com/TAULargeScaleWorkshop/RLAD/services/common/common-client"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CacheServiceClient struct {
	client.ServiceClientBase[service.CacheServiceClient]
}

func NewCacheServiceClient(addresses []string, service_name string) *CacheServiceClient {
	return &CacheServiceClient{
		ServiceClientBase: client.ServiceClientBase[service.CacheServiceClient]{
			RegistryAddresses: addresses,
			ServiceName:       service_name,
			CreateClient:      service.NewCacheServiceClient},
	}
}

func (obj *CacheServiceClient) Set(key string, value string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Set RPC function
	_, err = c.Set(context.Background(), &service.SetKeyValueReq{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("could not call Set: %v", err)
	}
	return nil
}

func (obj *CacheServiceClient) Get(key string) (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Get RPC function
	r, err := c.Get(context.Background(), &service.GetKeyReq{Key: key})
	if err != nil {
		return "", fmt.Errorf("could not call Get: %v", err)
	}
	return r.Value, nil
}

func (obj *CacheServiceClient) Delete(key string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Delete RPC function
	_, err = c.Delete(context.Background(), &service.GetKeyReq{Key: key})
	if err != nil {
		return fmt.Errorf("could not call Delete: %v", err)
	}
	return nil
}

func (obj *CacheServiceClient) IsAlive() (bool, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the IsAlive RPC function
	r, err := c.IsAlive(context.Background(), &emptypb.Empty{})
	if err != nil {
		return false, fmt.Errorf("could not call IsAlive: %v", err)
	}
	return r.Value, nil
}
