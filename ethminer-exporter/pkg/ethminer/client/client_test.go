package client

import(
	"testing"
)

func TestSendMessage(t *testing.T) {
	var api Api

	api = &Client{
		Host: "localhost",
		Port: 6666,
		ConnPoolSize: 2,
	}

	api.Init()

	api.SendMessage()

	api.Close()
}
