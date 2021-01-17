package client

import (
	"context"
	"fmt"
	"net"
	"os"

	retry "github.com/avast/retry-go"
	"github.com/sourcegraph/jsonrpc2"
)

type Api interface {
	Init()
	GetAddr() string
	// SendMessage()
	Ping() (string, error)
	// GetDetailStats() *jsonResponse
	Close()
}

type Client struct {
	Host         string
	Port         int
	ConnPoolSize int
	conns        chan *net.TCPConn
	netAddr      *net.TCPAddr
	addr         *string
}

type jsonApiObject struct {
	ID      int    `json:"id"`
	JsonRPC string `json:"jsonrpc"`
}

type jsonRequest struct {
	jsonApiObject `json:""`
	Method        string `json:"method"`
}

type jsonResponse struct {
	jsonApiObject `json:""`
	Result        *jsonResult `json:"result"`
}

type jsonResult map[string]interface{}

func (c *Client) GetDetailedStats() (*jsonResponse, error) {
	return nil, nil
}

func (c *Client) GetAddr() string {
	if c.addr != nil {
		return *c.addr
	}
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	c.addr = &addr
	return *c.addr
}

func (c *Client) returnConn(conn *net.TCPConn) {
	c.conns <- conn
}

type NullHandler struct{}

func (NullHandler) Handle(_ context.Context, _ *jsonrpc2.Conn, _ *jsonrpc2.Request) {}

func (c *Client) Ping() (string, error) {
	conn := <-c.conns
	//defer c.returnConn(conn)
	ctx := context.TODO()
	jsonRpcConn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, jsonrpc2.VarintObjectCodec{}), NullHandler{})

	var pong string
	pongErr := retry.Do(func() error {
		return jsonRpcConn.Call(ctx, "miner_ping", nil, &pong)
	})
	if pongErr != nil {
		return "", pongErr
	}

	return pong, nil
}

func (c *Client) Init() {
	var resolveErr error
	c.netAddr, resolveErr = net.ResolveTCPAddr("tcp", c.GetAddr())
	if resolveErr != nil {
		fmt.Println("ResolveTCPAddr failed:", resolveErr.Error())
		os.Exit(1)
	}

	c.conns = make(chan *net.TCPConn, c.ConnPoolSize)
	for i := 0; i < c.ConnPoolSize; i++ {
		conn, connErr := net.DialTCP("tcp", nil, c.netAddr)
		if connErr != nil {
			fmt.Println("ResolveTCPAddr failed:", connErr.Error())
			os.Exit(1)
		}
		c.conns <- conn
	}
}

func (c *Client) Close() {
	for i := 0; i < c.ConnPoolSize; i++ {
		conn := <-c.conns
		conn.Close()
	}
}
