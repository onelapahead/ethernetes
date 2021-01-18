package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/sourcegraph/jsonrpc2"
)

type loggingHandler struct{}

func (loggingHandler) Handle(ctx context.Context, _ *jsonrpc2.Conn, req *jsonrpc2.Request) {
	fmt.Println("Received jsonrpc method: ", req.Method)
}

type crlfObjectCodec struct{}

// WriteObject writes a JSON-RPC 2.0 object to the stream.
func (crlfObjectCodec) WriteObject(stream io.Writer, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	data = append(data, byte('\n'))

	if _, err := stream.Write(data); err != nil {
		return err
	}

	return nil
}

// ReadObject reads the next JSON-RPC 2.0 object from the stream
// and stores it in the value pointed to by v.
func (crlfObjectCodec) ReadObject(stream *bufio.Reader, v interface{}) error {
	responseBytes, readErr := stream.ReadBytes('\n')
	if readErr != nil {
		return readErr
	}

	unmarshalErr := json.Unmarshal(responseBytes, v)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

type jsonResult map[string]interface{}
