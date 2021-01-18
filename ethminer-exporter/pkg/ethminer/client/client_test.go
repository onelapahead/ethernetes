package client

import (
	"testing"
)

func TestPing(t *testing.T) {
	api := &ApiClient{
		Host:         "localhost",
		Port:         3333,
		ConnPoolSize: 2,
	}

	api.Init()
	defer api.Close()

	pong, err := api.Ping()
	if err != nil {
		t.Fail()
	}
	if pong != "pong" {
		t.Fail()
	}
}
