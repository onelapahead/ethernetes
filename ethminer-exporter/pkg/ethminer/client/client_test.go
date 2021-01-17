package client

import (
	"testing"
)

func TestPing(t *testing.T) {
	var api Api

	api = &Client{
		Host:         "localhost",
		Port:         6666,
		ConnPoolSize: 2,
	}

	api.Init()

	pong, err := api.Ping()
	if err != nil {
		t.Fail()
	}
	if pong != "pong" {
		t.Fail()
	}

	api.Close()
}
