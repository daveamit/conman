package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

func TestMain(m *testing.M) {
	defineConfigFlags()
	setLog(stringToLogLevel("debug"), "json")

	ctx, cancel := context.WithCancel(context.Background())
	go Serve(ctx)

	<-time.After(1 * time.Second)
	code := m.Run()
	cancel()
	os.Exit(code)
}

var client = resty.New().SetHostURL("http://localhost:8000")

func TestGetAllUsers(t *testing.T) {
	_, err := client.R().
		EnableTrace().
		Get("users")
	if err != nil {
		t.Error("Failed to retrive users", err)
	}
}

func TestGetUser(t *testing.T) {
	_, err := client.R().
		EnableTrace().
		Get("users/dummy-user")
	if err != nil {
		t.Error("Failed to retrive users", err)
	}
}
