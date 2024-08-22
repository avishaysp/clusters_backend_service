package RegServiceServant

import (
	"context"
	"fmt"
	"strings"
	"time"

	cacheservicecommon "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/common"
	testservicecommon "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// -------------------- Chord DHT encoding ---------------
// maps protocol to address
type NodeAddresses map[string]string

// encodes multiple protocols for the same node
// {"GRPC":"xxx", "MQ": "xxx"} -> "$GRPC$xxx$MQ$xxx"
func encodeProtocols(node_addresses NodeAddresses) string {
	enc_address := ""
	for key, value := range node_addresses {
		utils.Logger.Printf("")
		enc_address += "$" + key + "$" + value
	}
	return enc_address
}

// encodes a list of strings to a single string
// ["$GRPC$xxx$MQ$xxx", "$GRPC$yyy$MQ$yyy"] -> "$GRPC$xxx$MQ$xxx,$GRPC$yyy$MQ$yyy"
func encodeStrings(lst []string) string {
	return strings.Join(lst, ",")
}

// decodes multiple protocols for the same node
// "$GRPC$xxx$MQ$xxx" -> [("GRPC", "xxx"), ("MQ", "xxx")]
func decodeProtocols(enc string) NodeAddresses {
	lst := strings.Split(enc, "$")
	node_addresses := NodeAddresses{}
	// we skip i = 0 due to empty string after split
	for i := 1; i < len(lst); i += 2 {
		node_addresses[lst[i]] = lst[i+1]
	}
	return node_addresses
}

// decodes a single string of nodes to a list of strings
// "$GRPC$xxx$MQ$xxx,$GRPC$yyy$MQ$yyy" -> ["$GRPC$xxx$MQ$xxx", "$GRPC$yyy$MQ$yyy"]
func decodeStrings(enc string) []string {
	return strings.Split(enc, ",")
}

// -------------------- TestService --------------------
// TestServiceClient code (we duplicate some code for the IsAlive grpc/mq connections)
type TestServiceClient struct {
	// we have a specified address, not using the registry
	AddressGRPC  string
	AddressMQ    string
	CreateClient func(grpc.ClientConnInterface) testservicecommon.TestServiceClient
}

type CacheServiceClient struct {
	Address      string // we have a specified address, not using the registry
	CreateClient func(grpc.ClientConnInterface) cacheservicecommon.CacheServiceClient
}

func NewTestServiceClient(address_grpc string, address_mq string) *TestServiceClient {
	return &TestServiceClient{
		AddressGRPC:  address_grpc,
		AddressMQ:    address_mq,
		CreateClient: testservicecommon.NewTestServiceClient,
	}
}

func NewCacheServiceClient(address string) *CacheServiceClient {
	return &CacheServiceClient{
		Address:      address,
		CreateClient: cacheservicecommon.NewCacheServiceClient,
	}
}

func (obj *TestServiceClient) Connect() (res testservicecommon.TestServiceClient, closeFunc func(), err error) {
	// Set a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, obj.AddressGRPC, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty testservicecommon.TestServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.AddressGRPC, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
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

func (obj *CacheServiceClient) Connect() (res cacheservicecommon.CacheServiceClient, closeFunc func(), err error) {
	// Set a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, obj.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty cacheservicecommon.CacheServiceClient
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.Address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
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
