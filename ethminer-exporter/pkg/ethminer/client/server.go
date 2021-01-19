package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

type Server interface {
	Ping() (string, error)
	GetStatDetail() (*jsonResult, error)
}

type serverHandler struct {
	server Server
}

// Handle is called to handle a request. No other requests are handled
// until it returns. If you do not require strict ordering behavior
// of received RPCs, it is suggested to wrap your handler in
// AsyncHandler.
func (h serverHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// TODO default and replywitherror
	switch req.Method {
	case "miner_ping":
		pong, _ := h.server.Ping()
		_ = conn.Reply(ctx, req.ID, &pong)
	case "miner_getstatdetail":
		jsonResult, _ := h.server.GetStatDetail()
		_ = conn.Reply(ctx, req.ID, jsonResult)
	}
}

type MockServer struct {
	Host string
	Port int

	conns    []*jsonrpc2.Conn
	netAddr  *net.TCPAddr
	listener *net.TCPListener
	addr     string
	codec    jsonrpc2.ObjectCodec
	handler  jsonrpc2.Handler
}

func (s *MockServer) GetAddr() string {
	if s.addr != "" {
		return s.addr
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *MockServer) Init(ctx context.Context) {
	debug = true
	var resolveErr error
	s.codec = crlfObjectCodec{}
	s.handler = serverHandler{server: s}
	s.netAddr, resolveErr = net.ResolveTCPAddr("tcp", s.GetAddr())
	if resolveErr != nil {
		fmt.Println("ResolveTCPAddr failed:", resolveErr.Error())
		os.Exit(1)
	}

	var listenErr error
	s.listener, listenErr = net.ListenTCP("tcp", s.netAddr)
	if listenErr != nil {
		_ = s.listener.Close()
		fmt.Println("Failed to listen", listenErr.Error())
		os.Exit(1)
	}

	s.conns = make([]*jsonrpc2.Conn, 0)

	go func() {
		for true {
			ctx := context.TODO()
			conn, acceptErr := s.listener.AcceptTCP()
			if acceptErr != nil {
				time.Sleep(time.Second)
				continue
			}

			s.conns = append(s.conns, jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(conn, s.codec), s.handler))
		}
	}()
}

func (s *MockServer) Close() {
	_ = s.listener.Close()
	for _, conn := range s.conns {
		_ = conn.Close()
	}
}

func (s *MockServer) Ping() (string, error) {
	return "pong", nil
}

func (s *MockServer) GetStatDetail() (*jsonResult, error) {
	mockJsonResponse := `{"connection":{"connected":true,"switches":1,"uri":"stratum+tls12://0xf0bEA86827AE84B7a712a4Bc716a15C465be3878.rdu-01a@us1.ethermine.org:5555"},"devices":[{"_index":0,"_mode":"CUDA","hardware":{"name":"GeForce GTX 1080 7.93 GB","pci":"26:00.0","sensors":[0,0,0],"type":"GPU"},"mining":{"hashrate":"0x01426a98","pause_reason":null,"paused":false,"segment":["0xaa9ed703c60fa6bf","0xaa9ed704c60fa6bf"],"shares":[232,1,0,440]}}],"host":{"name":"4ac7c6e0e1e2","runtime":37163,"version":"ethminer-0.19.0-alpha.0"},"mining":{"difficulty":3999938964,"epoch":389,"epoch_changes":1,"hashrate":"0x01426a98","shares":[232,1,0,440]},"monitors":null}`
	result := &jsonResult{}

	err := json.Unmarshal([]byte(mockJsonResponse), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
