package client

import (
	"context"
	"fmt"
	"net"
	"os"

	retry "github.com/avast/retry-go"
	"github.com/sourcegraph/jsonrpc2"
)

type NetworkedClient interface {
	Init()
	GetAddr() string
	Close()
}

type Api interface {
	Ping() (string, error)
	// GetDetailStats() *jsonResponse
}

type ApiClient struct {
	Host         string
	Port         int
	ConnPoolSize int
	conns        chan *jsonrpc2.Conn
	netAddr      *net.TCPAddr
	addr         *string
}

type jsonResult map[string]interface{}

func (c *ApiClient) GetDetailedStats() (*jsonResult, error) {
	return nil, nil
}

func (c *ApiClient) GetAddr() string {
	if c.addr != nil {
		return *c.addr
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.addr = &addr
	return *c.addr
}

func (c *ApiClient) returnConn(conn *jsonrpc2.Conn) {
	c.conns <- conn
}

func (c *ApiClient) Ping() (string, error) {
	conn := <-c.conns
	defer c.returnConn(conn)
	ctx := context.TODO()

	var pong string
	pongErr := retry.Do(func() error {
		return conn.Call(ctx, "miner_ping", nil, &pong)
	})
	if pongErr != nil {
		return "", pongErr
	}

	return pong, nil
}

func (c *ApiClient) Init() {
	var resolveErr error
	ctx := context.TODO()
	codec := crlfObjectCodec{}
	requestHandler := loggingHandler{}
	c.netAddr, resolveErr = net.ResolveTCPAddr("tcp", c.GetAddr())
	if resolveErr != nil {
		fmt.Println("ResolveTCPAddr failed:", resolveErr.Error())
		os.Exit(1)
	}

	c.conns = make(chan *jsonrpc2.Conn, c.ConnPoolSize)
	for i := 0; i < c.ConnPoolSize; i++ {
		conn, connErr := net.DialTCP("tcp", nil, c.netAddr)
		if connErr != nil {
			fmt.Println("ResolveTCPAddr failed:", connErr.Error())
			os.Exit(1)
		}

		c.conns <- jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, codec), requestHandler)
	}
}

func (c *ApiClient) Close() {
	for i := 0; i < c.ConnPoolSize; i++ {
		conn := <-c.conns
		conn.Close()
	}
}
