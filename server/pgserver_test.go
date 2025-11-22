package server

import (
	"net"
	"os"
	"testing"

	"github.com/google/uuid"
	mconn "github.com/jordwest/mock-conn"
)

func TestQueryHandler(t *testing.T) {
	var (
		MockConn               = mconn.NewConn()
		started  chan struct{} = make(chan struct{})
	)
	defer MockConn.Close()

	send := make(chan net.Conn, 1)
	send <- MockConn.Server

	config := &Config{
		ID:         uuid.New(),
		RunStartUp: false,
		BackendKey: BackendKeyData{
			ProcessID: 123,
			SecretKey: 123,
		},
		Logger:  DefaultLogger(),
		Handler: &TestHandlerStruct{},
		InternalBuffer: &InternalBuffer{
			WriteBufferSize: 1024,
			Pool:            true,
		},
	}
	go _StartServer(t, send, started, config)
	<-started
	StartClient(t, MockConn.Client, map[string]Func{

		"Query": SendQuery,
	})
}

func TestParseHandler(t *testing.T) {
	var (
		MockConn               = mconn.NewConn()
		started  chan struct{} = make(chan struct{})
	)
	defer MockConn.Close()

	send := make(chan net.Conn, 1)
	send <- MockConn.Server

	config := &Config{
		ID:         uuid.New(),
		RunStartUp: false,
		BackendKey: BackendKeyData{
			ProcessID: 123,
			SecretKey: 123,
		},
		Logger:  DefaultLogger(),
		Handler: &TestHandlerStruct{},
		InternalBuffer: &InternalBuffer{
			WriteBufferSize: 1024,
			Pool:            false,
		},
	}
	go _StartServer(t, send, started, config)
	<-started
	StartClient(t, MockConn.Client, map[string]Func{
		"Parse": SendParse,
	})
}

func TestBindHandler(t *testing.T) {
	t.Parallel()

	var (
		MockConn               = mconn.NewConn()
		started  chan struct{} = make(chan struct{})
	)
	defer MockConn.Close()

	send := make(chan net.Conn, 1)
	send <- MockConn.Server

	config := &Config{
		ID:         uuid.New(),
		RunStartUp: false,
		BackendKey: BackendKeyData{
			ProcessID: 123,
			SecretKey: 123,
		},
		Logger:  DefaultLogger(),
		Handler: &TestHandlerStruct{},
		InternalBuffer: &InternalBuffer{
			WriteBufferSize: 1024,
		},
	}
	go _StartServer(t, send, started, config)
	<-started

	StartClient(t, MockConn.Client, map[string]Func{
		"Bind": SendBind,
	})
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
