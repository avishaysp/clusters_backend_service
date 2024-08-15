package common

import (
	"fmt"
	"math/rand"

	"github.com/pebbe/zmq4"

	common "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	RegServiceClient "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/client"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ServiceClientBase[client_t any] struct {
	RegistryAddresses []string // to connect to the registry service
	ServiceName       string   // to disover our server's nodes
	CreateClient      func(grpc.ClientConnInterface) client_t
}

func (obj *ServiceClientBase[client_t]) getMQNodes() ([]string, error) {
	regClient := RegServiceClient.NewRegServiceClient(obj.RegistryAddresses)
	nodes, err := regClient.Discover(obj.ServiceName + "MQ")
	if err != nil {
		Logger.Printf("getMQNodes(): Error calling Discover: %s", err)
	}
	return nodes, err
}

func (obj *ServiceClientBase[client_t]) ConnectMQ() (socket *zmq4.Socket, err error) {
	nodes, err := obj.getMQNodes()
	Logger.Printf("ConnectMQ(): got MQ nodes %v", nodes)
	if err != nil {
		Logger.Printf("Error calling getMQNodes: %s", err)
	}
	socket, err = zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	}
	for _, node := range nodes {
		Logger.Printf("ConnectMQ(): calling Connect on node %s", node)
		// TODO: consider adding timeout?
		err = socket.Connect(node)
		if err != nil {
			Logger.Printf("Failed to connect a zmq socket: %v", err)
		}
	}
	return socket, err
}

// randomly picks a service node address to connect to
func (obj *ServiceClientBase[client_t]) pickNode() string {
	regClient := RegServiceClient.NewRegServiceClient(obj.RegistryAddresses)
	nodes, err := regClient.Discover(obj.ServiceName)
	if err != nil || len(nodes) == 0 {
		return ""
	}
	randomIndex := rand.Intn(len(nodes))
	return nodes[randomIndex]
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
	// pick some node of the service
	node_address := obj.pickNode()
	if node_address == "" {
		var empty client_t
		return empty, nil, fmt.Errorf("no available nodes found")
	}
	Logger.Printf("Connect(): Got node address %s", node_address)

	// connect
	conn, err := grpc.Dial(node_address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		var empty client_t
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", node_address, err)
	}
	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}

func NewMarshaledCallParameter(method string, proto_data proto.Message) ([]byte, error) {
	var msg []byte

	// handle data
	data, err := proto.Marshal(proto_data)
	if err != nil {
		Logger.Printf("NewMarshaledCallParameter(): Marshal(proto_data) failed: %v\n", err)
	}

	// handle call params
	callParams := &common.CallParameters{Method: method, Data: data}
	msg, err = proto.Marshal(callParams)
	if err != nil {
		Logger.Printf("NewMarshaledCallParameter(): Marshal(callParams) failed: %v\n", err)
	}
	return msg, err
}
