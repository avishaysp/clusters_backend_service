package TestServiceClient

import (
	//context "context"
	"errors"
	"fmt"

	common "github.com/TAULargeScaleWorkshop/RLAD/services/common"
	client "github.com/TAULargeScaleWorkshop/RLAD/services/common/common-client"
	service "github.com/TAULargeScaleWorkshop/RLAD/services/test-service/common"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (obj *TestServiceClient) HelloWorldAsync() (func() (string, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("HelloWorldAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("HelloWorld", &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("HelloWorldAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("HelloWorldAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() (string, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return "", fmt.Errorf("HelloWorldAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return "", fmt.Errorf("HelloWorldAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		str := &wrapperspb.StringValue{}
		err = proto.Unmarshal(returnValue.Data, str)
		if err != nil {
			return "", fmt.Errorf("HelloWorldAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return str.Value, err
	}

	return ret, nil
}

func (obj *TestServiceClient) HelloToUserAsync(userName string) (func() (string, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("HelloToUserAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("HelloToUser", wrapperspb.String(userName))
	if err != nil {
		return nil, fmt.Errorf("HelloToUserAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("HelloToUserAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() (string, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return "", fmt.Errorf("HelloToUserAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return "", fmt.Errorf("HelloToUserAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		str := &wrapperspb.StringValue{}
		err = proto.Unmarshal(returnValue.Data, str)
		if err != nil {
			return "", fmt.Errorf("HelloToUserAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return str.Value, err
	}

	return ret, nil
}

func (obj *TestServiceClient) StoreAsync(key string, value string) (func() error, error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("StoreAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("Store", &service.StoreKeyValue{Key: key, Value: value})
	if err != nil {
		return nil, fmt.Errorf("StoreAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("StoreAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() error {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return fmt.Errorf("StoreAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return fmt.Errorf("StoreAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle only error (returnValue.Data is empty - Store is a void function)
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return err
	}

	return ret, nil
}

func (obj *TestServiceClient) GetAsync(key string) (func() (string, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("GetAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("Get", wrapperspb.String(key))
	if err != nil {
		return nil, fmt.Errorf("GetAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("GetAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() (string, error) {
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return "", fmt.Errorf("GetAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return "", fmt.Errorf("GetAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		value := &wrapperspb.StringValue{}
		err = proto.Unmarshal(returnValue.Data, value)
		if err != nil {
			return "", fmt.Errorf("GetAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return value.Value, err
	}

	return ret, nil
}

func (obj *TestServiceClient) IsAliveAsync() (func() (bool, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("IsAliveAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("IsAlive", &emptypb.Empty{})
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

func (obj *TestServiceClient) ExtractLinksFromURLAsync(url string, depth int32) (func() ([]string, error), error) {
	mqsocket, err := obj.ConnectMQ()
	if err != nil {
		return nil, fmt.Errorf("ExtractLinksFromURLAsync(): ConnectMQ failed: %v\n", err)
	}

	// packing
	msg, err := client.NewMarshaledCallParameter("ExtractLinksFromURL", &service.ExtractLinksFromURLParameters{Url: url, Depth: depth})
	if err != nil {
		return nil, fmt.Errorf("ExtractLinksFromURLAsync(): NewMarshaledCallParameter failed: %v\n", err)
	}

	_, err = mqsocket.SendBytes(msg, 0)
	if err != nil {
		return nil, fmt.Errorf("ExtractLinksFromURLAsync(): SendBytes failed: %v\n", err)
	}

	// return function (future pattern)
	ret := func() ([]string, error) {
		var links []string
		defer mqsocket.Close()
		rv, err := mqsocket.RecvBytes(0)
		if err != nil {
			return links, fmt.Errorf("ExtractLinksFromURLAsync(): RecvBytes failed: %v\n", err)
		}

		returnValue := &common.ReturnValue{}
		// handle return value
		err = proto.Unmarshal(rv, returnValue)
		if err != nil {
			return links, fmt.Errorf("ExtractLinksFromURLAsync(): Unmarshal(rv) failed: %v\n", err)
		}

		// handle data
		value := &service.ExtractLinksFromURLReturnedValue{}
		err = proto.Unmarshal(returnValue.Data, value)
		if err != nil {
			return links, fmt.Errorf("ExtractLinksFromURLAsync(): Unmarshal(returnValue.Data) failed: %v\n", err)
		}
		links = value.Links

		// error
		if returnValue.Error != "" {
			err = errors.New(returnValue.Error)
		}
		return links, err
	}

	return ret, nil
}
