package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgproto3"
)

type _common interface {
	// AuthMethod
	ID() uuid.UUID
	Context() context.Context
	Write([]byte) (int, error)
	SetClientProtocolVer(uint32)
}

type Context interface {
	_common
	WriteMsg(pgproto3.BackendMessage) error
	ReadyForQuery(TxStatus byte) error
}

type StartUpCtx interface {
	_common
	Read([]byte) (int, error)
}

// _common Interface

func (s *ctx_server) ID() uuid.UUID {
	return s.id
}

func (s *ctx_server) Context() context.Context {
	return s.ctx
}

func (s *ctx_server) Write(buf []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, err = s.conn.Write(buf)
	return
}

func (s *ctx_server) SetClientProtocolVer(v uint32) {

}

// Context Interface
func (s *ctx_server) WriteMsg(msg pgproto3.BackendMessage) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.wbuf.Data) == 0 {
		s.wbuf.Data, err = msg.Encode(nil)
		return
	}
	s.wbuf.Data, err = msg.Encode(s.wbuf.Data)
	return
}

func (s *ctx_server) ReadyForQuery(TxStatus byte) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.wbuf.Data, err = (&pgproto3.ReadyForQuery{TxStatus: TxStatus}).Encode(s.wbuf.Data)
	return err
}

// StartUpCtx Interface
func (s *ctx_server) Read(b []byte) (int, error) {
	return s.conn.Read(b)
}
