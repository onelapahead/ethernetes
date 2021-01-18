package client

import (
	"fmt"
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
		t.Fatal(err)
	}
	if pong != "pong" {
		t.Fail()
	}

	statsResult, err := api.GetDetailedStats()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(statsResult)
}
