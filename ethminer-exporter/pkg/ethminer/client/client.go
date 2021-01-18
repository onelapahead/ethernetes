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
	getAddr() string
	Close()
}

type Api interface {
	Ping() (string, error)
	GetDetailStats() (*jsonResult, error)
}

type ApiClient struct {
	Host         string
	Port         int
	ConnPoolSize int

	conns   chan *jsonrpc2.Conn
	netAddr *net.TCPAddr
	addr    *string
	codec   jsonrpc2.ObjectCodec
	handler jsonrpc2.Handler
}

func (c *ApiClient) returnConn(conn *jsonrpc2.Conn) {
	c.conns <- conn
}

func (c *ApiClient) newConn(ctx context.Context) (*jsonrpc2.Conn, error) {
	conn, connErr := net.DialTCP("tcp", nil, c.netAddr)
	if connErr != nil {
		return nil, connErr
	}

	return jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, c.codec), c.handler), nil
}

func (c *ApiClient) getAddr() string {
	if c.addr != nil {
		return *c.addr
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.addr = &addr
	return *c.addr
}

func (c *ApiClient) Init() {
	var resolveErr error
	ctx := context.TODO()
	c.codec = crlfObjectCodec{}
	c.handler = loggingHandler{}
	c.netAddr, resolveErr = net.ResolveTCPAddr("tcp", c.getAddr())
	if resolveErr != nil {
		fmt.Println("ResolveTCPAddr failed:", resolveErr.Error())
		os.Exit(1)
	}

	c.conns = make(chan *jsonrpc2.Conn, c.ConnPoolSize)
	for i := 0; i < c.ConnPoolSize; i++ {
		conn, connErr := c.newConn(ctx)
		if connErr != nil {
			fmt.Println("ResolveTCPAddr failed:", connErr.Error())
			os.Exit(1)
		}

		c.conns <- conn
	}
}

func (c *ApiClient) Close() {
	for i := 0; i < c.ConnPoolSize; i++ {
		conn := <-c.conns
		conn.Close()
	}
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
		conn.Close()
		conn, _ = c.newConn(ctx)

		return "", pongErr
	}

	return pong, nil
}

func (c *ApiClient) GetDetailedStats() (*jsonResult, error) {
	fmt.Println("blah")
	conn := <-c.conns
	defer c.returnConn(conn)
	ctx := context.TODO()

	var result *jsonResult
	err := retry.Do(func() error {
		return conn.Call(ctx, "miner_getstatdetail", nil, result)
	})
	if err != nil {
		conn.Close()
		conn, _ = c.newConn(ctx)

		return nil, err
	}

	return result, nil
}
