package client

import (
	"context"
	"fmt"
	"testing"
)

func TestPing(t *testing.T) {
	api := &ApiClient{
		Host:         "localhost",
		Port:         3333,
		ConnPoolSize: 2,
	}

	server := &MockServer{
		Host: "localhost",
		Port: 3333,
	}
	ctx := context.TODO()

	server.Init(ctx)
	api.Init(ctx)
	defer api.Close()
	defer server.Close()

	pong, err := api.Ping(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if pong != "pong" {
		t.Fail()
	}

	statsResult, err := api.GetStatDetail(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(statsResult)
}
