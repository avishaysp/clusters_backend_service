package RegServiceServant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pebbe/zmq4"

	cacheservicecommon "github.com/TAULargeScaleWorkshop/RLAD/services/cache-service/common"
	"github.com/TAULargeScaleWorkshop/RLAD/services/common"
	testservicecommon "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"
	"github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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

// -------------------- TestServiceMQ --------------------
func (obj *TestServiceClient) ConnectMQ() (socket *zmq4.Socket, err error) {
	socket, err = zmq4.NewSocket(zmq4.REQ)
	utils.Logger.Printf("ConnectMQ(): created NewSocket %v", socket)
	if err != nil {
		utils.Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	}
	utils.Logger.Printf("ConnectMQ(): calling Connect for address: %s", obj.AddressMQ)
	err = socket.Connect(obj.AddressMQ)
	if err != nil {
		utils.Logger.Printf("Failed to connect a zmq socket: %v", err)
	}
	return socket, err
}

func (obj *TestServiceClient) NewMarshaledCallParameter(method string, proto_data proto.Message) ([]byte, error) {
	var msg []byte

	// handle data
	data, err := proto.Marshal(proto_data)
	if err != nil {
		utils.Logger.Printf("NewMarshaledCallParameter(): Marshal(proto_data) failed: %v\n", err)
	}

	// handle call params
	callParams := &common.CallParameters{Method: method, Data: data}
	msg, err = proto.Marshal(callParams)
	if err != nil {
		utils.Logger.Printf("NewMarshaledCallParameter(): Marshal(callParams) failed: %v\n", err)
	}
	return msg, err
}

func (obj *TestServiceClient) IsAliveAsync() (func() (bool, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := obj.NewMarshaledCallParameter("IsAlive", &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() (bool, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		r := &wrapperspb.BoolValue{}
		err = proto.Unmarshal(returnValue.Data, r)
		if err != nil {
			return false, fmt.Errorf("IsAliveAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return r.Value, err
	}

	return ret, nil
}
