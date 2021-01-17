package client

import (
	"fmt"
	"net"
	"os"

	retry "github.com/avast/retry-go"
)

type Api interface {
	Init()
	GetAddr() string
	SendMessage()
	// Ping() *jsonResponse
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
	jsonApiObject
	Method string `json:"method"`
}

type jsonResponse struct {
	jsonApiObject
	Result *jsonResult `json:"result"`
}

type jsonResult map[string]interface{}

func (c *Client) Ping() *jsonResponse {

	return nil
}

func (c *Client) GetDetailedStats() *jsonResponse {

	return nil
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

func (c *Client) SendMessage() {
	conn := <-c.conns
	defer c.returnConn(conn)

	retry.Do(func() error {

		_, sendErr := conn.Write([]byte("hello!\n"))
		if sendErr != nil {
			return sendErr
		}

		fmt.Println("sent hello!")
		return nil
	})

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
