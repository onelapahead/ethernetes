package client

import (
	"fmt"
	"net"
	"os"
)

type Api interface {
	Init()
	GetAddr() string
	SendMessage()
	Close()
}

type Client struct {
	Host string
	Port int
	ConnPoolSize int
	conns chan *net.TCPConn
	netAddr *net.TCPAddr
	addr *string
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

	_, sendErr := conn.Write([]byte("hello!"))
	if sendErr != nil {

	}

	fmt.Println("sent hello!")
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
