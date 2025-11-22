package server

import (
	"context"
	"net"
	"testing"
	"time"
)

func _StartServer(t *testing.T, conn_c chan net.Conn, started chan struct{}, config *Config) {
	// listener, err := net.Listen("tcp", _host_addr)
	// defer func() {
	// 	err := listener.Close()
	// 	assert.Nil(t, err)
	// }()
	// assert.Nil(t, err)
	// addr <- listener.Addr().String()

	started <- struct{}{}
	for {
		conn := <-conn_c
		go handleConnection(t, conn, config)

	}
}

func handleConnection(t *testing.T, conn net.Conn, config *Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer conn.Close()

	sv, err := New(config)
	if err != nil {
		t.Errorf("Error Occured when creating Server instance:  %e", err)
	}
	sv.Accept(ctx, conn)
	// sv.Err
}
