package client

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/sourcegraph/jsonrpc2"
)

type Networked interface {
	Init(ctx context.Context)
	GetAddr() string
	Close()
}

type Api interface {
	Ping(ctx context.Context) (string, error)
	GetStatDetail(ctx context.Context) (*jsonResult, error)
}

type ApiClient struct {
	Host         string
	Port         int
	ConnPoolSize int

	conns   chan *jsonrpc2.Conn
	netAddr *net.TCPAddr
	addr    string
	codec   jsonrpc2.ObjectCodec
	handler jsonrpc2.Handler
}

func (c *ApiClient) returnConn(conn *jsonrpc2.Conn) {
	c.conns <- conn
}

func (c *ApiClient) returnConnAsync(conn *jsonrpc2.Conn) {
	go c.returnConn(conn)
}

func (c *ApiClient) newConn(ctx context.Context) (*jsonrpc2.Conn, error) {
	conn, connErr := net.DialTCP("tcp", nil, c.netAddr)
	if connErr != nil {
		return nil, connErr
	}

	return jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, c.codec), jsonrpc2.AsyncHandler(c.handler)), nil
}

func (c *ApiClient) GetAddr() string {
	if c.addr != "" {
		return c.addr
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ApiClient) Init(ctx context.Context) {
	var resolveErr error
	c.codec = crlfObjectCodec{}
	c.handler = loggingHandler{}
	c.netAddr, resolveErr = net.ResolveTCPAddr("tcp", c.GetAddr())
	if resolveErr != nil {
		fmt.Println("ResolveTCPAddr failed:", resolveErr.Error())
		os.Exit(1)
	}

	c.conns = make(chan *jsonrpc2.Conn, c.ConnPoolSize)
	for i := 0; i < c.ConnPoolSize; i++ {
		conn, connErr := c.newConn(ctx)
		if connErr != nil {
			fmt.Println("Failed to create connection:", connErr.Error())
			os.Exit(1)
		}

		c.conns <- conn
	}
}

func (c *ApiClient) Close() {
	for i := 0; i < c.ConnPoolSize; i++ {
		conn := <-c.conns
		_ = conn.Close()
	}
}

func (c *ApiClient) do(ctx context.Context, method string, obj interface{}, opts ...retry.Option) error {
	conn := <-c.conns
	defer c.returnConnAsync(conn)

	err := retry.Do(func() error {
		return conn.Call(ctx, method, nil, obj)
	}, opts...)
	if err != nil {
		// TODO if err is network related
		// best effort to close and create new conn
		// will still be returned to the pool async
		_ = conn.Close()
		conn, _ = c.newConn(ctx)

		return err
	}

	return nil
}

func (c *ApiClient) Ping(ctx context.Context) (string, error) {
	var pong string
	err := c.do(ctx, "miner_ping", &pong, retry.Attempts(1), retry.Delay(10*time.Millisecond))
	if err != nil {
		return "", err
	}

	return pong, nil
}

func (c *ApiClient) GetStatDetail(ctx context.Context) (*jsonResult, error) {
	result := &jsonResult{}
	err := c.do(ctx, "miner_getstatdetail", result, retry.Attempts(1), retry.Delay(400*time.Millisecond))
	if err != nil {
		return nil, err
	}

	return result, nil
}
