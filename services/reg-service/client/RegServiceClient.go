package RegServiceClient

import (
	context "context"
	"fmt"
	"math/rand"
	"strings"

	service "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/common"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Note: Code Duplication from RegServiceUtils (to prevent dht copy)
func decodeProtocols(enc string) map[string]string {
	lst := strings.Split(enc, "$")
	node_addresses := map[string]string{}
	// we skip i = 0 due to empty string after split
	for i := 1; i < len(lst); i += 2 {
		node_addresses[lst[i]] = lst[i+1]
	}
	return node_addresses
}

// Note: Code Duplication to prevent importing ServiceClientBase creating an import cycle
type RegServiceClient struct {
	// We don't need ServiceName field like in ServiceClientBase (since we are the registry)
	RegistryAddresses []string
	CreateClient      func(grpc.ClientConnInterface) service.RegServiceClient
}

func (obj *RegServiceClient) Connect() (res service.RegServiceClient, closeFunc func(), err error) {
	randomIndex := rand.Intn(len(obj.RegistryAddresses))
	conn, err := grpc.Dial(obj.RegistryAddresses[randomIndex], grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty service.RegServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.RegistryAddresses[0], err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}

func NewRegServiceClient(addresses []string) *RegServiceClient {
	return &RegServiceClient{
		RegistryAddresses: addresses, CreateClient: service.NewRegServiceClient,
	}
}

func convertFullAddress(node_addresses map[string]string) []*service.FullAddress {
	full_addresses := []*service.FullAddress{}
	for key, value := range node_addresses {
		fa := service.FullAddress{Protocol: key, Address: value}
		full_addresses = append(full_addresses, &fa)
	}
	return full_addresses
}

func (obj *RegServiceClient) Register(service_name string, addresses map[string]string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to registry server: %v", err)
	}
	defer closeFunc()

	full_addresses := convertFullAddress(addresses)

	// Call the Register RPC function
	_, err = c.Register(context.Background(), &service.UpdateRegistryParameters{
		ServiceName: service_name, Addresses: full_addresses})
	if err != nil {
		return fmt.Errorf("could not call Register: %v", err)
	}
	return nil
}

func (obj *RegServiceClient) Unregister(service_name string, addresses map[string]string) error {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to registry server: %v", err)
	}
	defer closeFunc()

	full_addresses := convertFullAddress(addresses)

	// Call the Unregister RPC function
	_, err = c.Unregister(context.Background(), &service.UpdateRegistryParameters{
		ServiceName: service_name, Addresses: full_addresses})
	if err != nil {
		return fmt.Errorf("could not call Unregister: %v", err)
	}
	return nil
}

func (obj *RegServiceClient) Discover(service_name string, protocol string) ([]string, error) {
	c, closeFunc, err := obj.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to registry server: %v", err)
	}
	defer closeFunc()

	// Call the Discover RPC function
	discovered, err := c.Discover(context.Background(), wrapperspb.String(service_name))
	if err != nil {
		return nil, fmt.Errorf("could not call Discover: %v", err)
	}

	// Return only the addresses for the given protocol
	addresses := []string{}

	for _, enc_address := range discovered.Addresses {
		node_addresses := decodeProtocols(enc_address)
		for key, value := range node_addresses {
			if key == protocol {
				addresses = append(addresses, value)
			}
		}
	}

	return addresses, nil
}
