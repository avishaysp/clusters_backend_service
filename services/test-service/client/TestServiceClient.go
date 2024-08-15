package TestServiceClient

import (
	context "context"
	"fmt"

	client "github.com/TAULargeScaleWorkshop/RLAD/services/common/common-client"
	service "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TestServiceClient struct {
	client.ServiceClientBase[service.TestServiceClient]
}

func NewTestServiceClient(addresses []string, service_name string) *TestServiceClient {
	return &TestServiceClient{
		ServiceClientBase: client.ServiceClientBase[service.TestServiceClient]{
			RegistryAddresses: addresses,
			ServiceName:       service_name,
			CreateClient:      service.NewTestServiceClient},
	}
}
func (obj *TestServiceClient) HelloWorld() (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloWorld(context.Background(), &emptypb.Empty{})
	if err != nil {
		return "", fmt.Errorf("could not call HelloWorld: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) HelloToUser(userName string) (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the HelloWorld RPC function
	r, err := c.HelloToUser(context.Background(), wrapperspb.String(userName))
	if err != nil {
		return "", fmt.Errorf("could not call HelloToUser: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) Store(key string, value string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Store RPC function
	_, err = c.Store(context.Background(), &service.StoreKeyValue{Key: key, Value: value})
	if err != nil {
		return fmt.Errorf("could not call Store: %v", err)
	}
	return nil
}

func (obj *TestServiceClient) Get(key string) (string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the Get RPC function
	r, err := c.Get(context.Background(), wrapperspb.String(key))
	if err != nil {
		return "", fmt.Errorf("could not call Get: %v", err)
	}
	return r.Value, nil
}

func (obj *TestServiceClient) WaitAndRand(seconds int32) (func() (int32, error), error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to server: %v", err)
	}
	r, err := c.WaitAndRand(context.Background(), wrapperspb.Int32(seconds))
	if err != nil {
		return nil, fmt.Errorf("could not call Get: %v", err)
	}
	res := func() (int32, error) {
		defer closeFunc()
		x, err := r.Recv()
		return x.Value, err
	}
	return res, nil
}

func (obj *TestServiceClient) IsAlive() (bool, error) {
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

func (obj *TestServiceClient) ExtractLinksFromURL(url string, depth int32) ([]string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to server: %v", err)
	}
	defer closeFunc()
	// Call the ExtractLinksFromURL RPC function
	r, err := c.ExtractLinksFromURL(context.Background(), &service.ExtractLinksFromURLParameters{Url: url, Depth: depth})
	if err != nil {
		return nil, fmt.Errorf("could not call ExtractLinksFromURL: %v", err)
	}
	return r.Links, nil
}
