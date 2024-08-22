package common

import (
	"fmt"
	"net"

	common "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	RegServiceClient "github.com/TAULargeScaleWorkshop/RLAD/services/reg-service/client"
	. "github.com/TAULargeScaleWorkshop/RLAD/utils"
	"github.com/pebbe/zmq4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// maps protocol to address (duplicate from RegServiceUtils to prevent dht copy)
type NodeAddresses map[string]string

func bindMQToService(listenPort int, messageHandler func(method string,
	parameters []byte) (response proto.Message, err error)) (startMQ func(), listeningAddress string) {
	Logger.Printf("bindMQToService() called with listenPort %d", listenPort)

	socket, err := zmq4.NewSocket(zmq4.REP)
	if err != nil {
		Logger.Fatalf("Failed to create a new zmq socket: %v", err)
	}
	if listenPort == 0 {
		listeningAddress = "tcp://127.0.0.1:*"
	} else {
		listeningAddress = fmt.Sprintf("tcp://127.0.0.1:%v", listenPort)
	}

	Logger.Printf("bindMQToService() calling Bind on %s", listeningAddress)
	err = socket.Bind(listeningAddress)
	if err != nil {
		Logger.Fatalf("Failed to bind a zmq socket: %v", err)
	}

	listeningAddress, err = socket.GetLastEndpoint()
	if err != nil {
		Logger.Fatalf("Failed to get listetning address of zmq socket: %v", err)
	}
	Logger.Printf("bindMQToService() GetLastEndpoint returned %s", listeningAddress)

	startMQ = func() {
		for {
			Logger.Printf("startMQ(): entered for loop")
			// TODO: flag for no wait
			data, readerr := socket.RecvBytes(0)
			Logger.Printf("startMQ(): called RecvBytes")
			if err != nil {
				Logger.Printf("Failed to receive bytes from MQ socket: %v\n", readerr)
				continue
			}
			if len(data) == 0 {
				Logger.Printf("No data to process")
				continue
			}
			Logger.Printf("data len: %v\n", len(data))

			// handle the request in a new goroutine
			go func() {
				Logger.Printf("startMQ func(): unpacking")
				// unpacking
				callParams := &common.CallParameters{}
				err := proto.Unmarshal(data, callParams)
				if err != nil {
					Logger.Printf("Unmarshal failed: %v\n", err)
					return
				}

				// call the method
				Logger.Printf("startMQ func(): calling messageHandler")
				response, err := messageHandler(callParams.Method, callParams.Data)

				// packing
				Logger.Printf("startMQ func(): got response")
				returnValue := &common.ReturnValue{}
				if err != nil {
					Logger.Printf("messageHandler failed with: %v\n", err)
					returnValue.Error = err.Error() // get the string for the error
				} else { // no error, get the data
					returnValue.Data, err = proto.Marshal(response)
					if err != nil {
						Logger.Printf("Marshal (1) failed: %v\n", err)
						returnValue.Error = err.Error()
					}
				}

				// process the full returnValue (data + err)
				responseData, err := proto.Marshal(returnValue)
				if err != nil {
					Logger.Printf("Marshal (2) failed: %v\n", err)
					return
				}

				_, senderr := socket.SendBytes(responseData, 0)
				if senderr != nil {
					Logger.Printf("SendBytes failed: %v\n", err)
					return
				}
			}()
		}
	}
	return startMQ, listeningAddress
}

func startgRPC(listenPort int) (listeningAddress string, grpcServer *grpc.Server, startListening func(), err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
	if err != nil {
		Logger.Fatalf("failed to listen: %v", err)
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

// regAddresses - the registry service addresses to connect to (for all the services who aren't the registry)
func Start(serviceName string, grpcListenPort int, regAddresses []string, bindgRPCToService func(s grpc.ServiceRegistrar),
	messageHandler func(method string, parameters []byte) (response proto.Message, err error)) (err error) {
	Logger.Printf("Start(%s, %d)\n", serviceName, grpcListenPort)

	// start the gRPC
	listeningAddress, grpcServer, startListening, err := startgRPC(grpcListenPort)
	if err != nil {
		return err
	}
	bindgRPCToService(grpcServer)

	// create a list of full addresses
	node_addresses := map[string]string{}

	// first, add the GRPC
	node_addresses["GRPC"] = listeningAddress

	// mq
	if messageHandler != nil {
		start_mq, listening_address_mq := bindMQToService(0, messageHandler)

		// add MQ
		node_addresses["MQ"] = listening_address_mq

		Logger.Printf("MQ: %s calling start_mq on %s\n", serviceName, listening_address_mq)
		go start_mq()
	}

	unregister := registerAddress(serviceName, regAddresses, node_addresses)
	defer unregister()

	Logger.Printf("GRPC: %s starts listening on %s\n", serviceName, listeningAddress)
	startListening()
	return nil
}

func registerAddress(serviceName string, regAddresses []string, node_addresses map[string]string) (unregister func()) {
	regClient := RegServiceClient.NewRegServiceClient(regAddresses)
	err := regClient.Register(serviceName, node_addresses)
	if err != nil {
		Logger.Fatalf("Failed to register to registry service: %v", err)
	}
	return func() {
		regClient.Unregister(serviceName, node_addresses)
	}
}
