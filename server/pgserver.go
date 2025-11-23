// BSD 3-Clause License

// Copyright (c) 2025, Patrick Innocent

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package server

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgproto3"
)

// TODO: FILE to *CONFIG using gob encoding
// func ParseConfig(byt []byte, enc_type ) (*Config, error) {
// 	return &Config{}, nil
// }q	az

type Config struct {
	// AuthMethod       AuthMethod
	ID               uuid.UUID
	Debug            bool
	Handler          Handler
	RunStartUp       bool
	BackendKey       BackendKeyData
	Logger           Logger
	LimitEOF         uint8
	PasswordAuthType uint32
	InternalBuffer   *InternalBuffer
}

type InternalBuffer struct {
	WriteBufferSize uint16
	Pool            bool
}

type server struct {
	id             uuid.UUID
	mu             sync.RWMutex
	ctx            context.Context
	once           sync.Once
	conn           net.Conn
	backend        *pgproto3.Backend
	h              Handler
	backendKeyData BackendKeyData
	log            Logger
	cancelRequest  chan struct{}
	debug          bool

	runStartUp        bool
	curentEOF         uint8
	limitEOF          uint8
	currenthandlertag HandlerTag
	passwordAuthType  uint32
	err               error
	// internal buffer
	wbuf            *Bytes
	writebuffersize uint16
	wbufpool        bool
	//
	//
	ctx_s *ctx_server
}

type ctx_server struct {
	*server
}

func New(cfg *Config) (*server, error) {
	if cfg.InternalBuffer.WriteBufferSize == 0 {
		cfg.InternalBuffer.WriteBufferSize = 1024 // Default buffer size
	}
	if cfg.LimitEOF == 0 {
		cfg.LimitEOF = 1
	}
	if cfg.Handler == nil {
		return nil, errors.New("parameter - handler can not be nil")
	}

	sv := &server{
		id:             cfg.ID,
		h:              cfg.Handler,
		backendKeyData: cfg.BackendKey,
		log:            cfg.Logger,
		cancelRequest:  make(chan struct{}, 1), // Buffered to prevent blocking
		// wbuf:              Bytes{Cap: cfg.InternalBuffer.WriteBufferSize},
		runStartUp:        cfg.RunStartUp,
		limitEOF:          cfg.LimitEOF,
		curentEOF:         0,
		writebuffersize:   cfg.InternalBuffer.WriteBufferSize,
		currenthandlertag: HandlerTagNone,
	}
	sv.ctx_s = &ctx_server{
		server: sv,
	}

	return sv, nil
}

// Handler = nil -> No Set
func (s *server) SetHandler(h Handler) {
	if h == nil {
		return
	}
	s.h = h
}

func (s *server) SetLogger(l Logger) {
	s.log = l
}

func (s *server) Accept(ctx context.Context, conn net.Conn) {
	var (
		msgChan chan pgproto3.FrontendMessage = make(chan pgproto3.FrontendMessage)
		errChan chan error                    = make(chan error, 1)
		err     error
	)
	s.ctx = ctx
	s.conn = conn

	// Handle StartUp Message
	if s.runStartUp {
		s.currenthandlertag = HandlerTagStartUpMsg

		err = s.h.HandlerStartupMessage(s.ctx_s)
		if err != nil {
			if (err == io.ErrUnexpectedEOF) || (err == io.EOF) {
				s.h.HandleError(&Error{
					ErrorHanlder: s.currenthandlertag,
					Err:          err,
					Level:        "FATAL",
					Comment:      "-connection was closed-",
				})
				return
			}
			s.h.HandleError(&Error{
				ErrorHanlder: s.currenthandlertag,
				Err:          err,
				Level:        "FTAL",
				Comment:      "-",
			})
			return
		}
	}

	s.backend = pgproto3.NewBackend(s.conn, s.conn)
	idx := func() uint16 {
		if s.wbufpool {
			return Count()
		}
		return BytesPure
	}()
	s.wbuf = &Bytes{
		Cap:   s.writebuffersize,
		Index: idx,
	}
	s.wbuf.Init()

	// Handle Receiving Message
	go func() {
		defer close(msgChan)
		for {
			msg, err := s.backend.Receive()
			if err != nil {
				errChan <- err
				return
			}
			select {
			case msgChan <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Main Loop
	for {

		select {
		case msg := <-msgChan:
			s.wbuf.Get()

			err = s.handleMsg(msg)
			if err != nil {
				s.h.HandleError(&Error{
					ErrorHanlder: s.currenthandlertag,
					Err:          err,
					Level:        "FATAL",
					Comment:      "- there was an error returned, passing the message to the Handler - ",
				})
				return
			}
			_, err = s.write()
			if err != nil {
				s.h.HandleError(&Error{
					ErrorHanlder: s.currenthandlertag,
					Err:          err,
					Level:        "FATAL",
					Comment:      "- write error -",
				})
				return
			}

		case err = <-errChan:
			// Handle Error
			if (err == io.ErrUnexpectedEOF) || (err == io.EOF) {
				s.h.HandleError(&Error{
					ErrorHanlder: HandlerTagEOF,
					Err:          err,
					Level:        "FATAL",
					Comment:      "-connection was closed-",
				})
			}

			s.h.HandleError(&Error{
				ErrorHanlder: s.currenthandlertag,
				Err:          err,
				Level:        "FATAL",
			})
			return

		case <-ctx.Done():
			// Handle Context Cancel - and return
			s.h.HandleError(&Error{
				ErrorHanlder: s.currenthandlertag,
				Err:          ctx.Err(),
				Level:        "FATAL",
				Comment:      "- context error -",
			})
			return

		case <-s.cancelRequest:
			s.currenthandlertag = HandlerTagCancelRequest
			s.h.HandlerCancelRequest(s.ctx_s, &pgproto3.CancelRequest{
				ProcessID: s.backendKeyData.ProcessID,
				SecretKey: s.backendKeyData.SecretKey,
			})
		}

	}

}
func (s *server) write() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.wbuf.Data) == 0 {
		return 0, nil
	}
	n, err := s.conn.Write(s.wbuf.Data)
	if err == nil {
		s.wbuf.Put()
	}
	return n, err
}

func (s *server) CancelRequest(pid uint32, skey uint32) {
	if s.backendKeyData.ProcessID == pid && s.backendKeyData.SecretKey == skey {
		select {
		case s.cancelRequest <- struct{}{}:
		default:
			// Channel is full, a cancel is already pending.
			// Log this or ignore it.
		}
	}
}

func (s *server) handleMsg(msg pgproto3.FrontendMessage) error {
	var err error

	switch _msg := msg.(type) {

	case *pgproto3.Query:
		s.currenthandlertag = HandlerTagQuery
		err = s.h.HandlerQuery(s.ctx_s, _msg)

	case *pgproto3.Parse:
		s.currenthandlertag = HandlerTagParse
		err = s.h.HandlerParse(s.ctx_s, _msg)

	case *pgproto3.Bind:
		s.currenthandlertag = HandlerTagBind
		err = s.h.HandlerBind(s.ctx_s, _msg)

	case *pgproto3.Execute:
		s.currenthandlertag = HandlerTagExcute
		err = s.h.HandlerExecute(s.ctx_s, _msg)

	case *pgproto3.Sync:
		s.currenthandlertag = HandlerTagSync
		err = s.h.HandlerSync(s.ctx_s, _msg)

	case *pgproto3.FunctionCall:
		s.currenthandlertag = HandlerTagFunctionCall
		err = s.h.HandlerFunctionCall(s.ctx_s, _msg)

	case *pgproto3.Close:
		s.currenthandlertag = HandlerTagClose
		err = s.h.HandlerClose(s.ctx_s, _msg)

	case *pgproto3.Terminate:
		s.currenthandlertag = HandlerTagTerminate
		err = s.h.HandlerTerminate(s.ctx_s, _msg)

	case *pgproto3.CancelRequest:
		if s.backendKeyData.ProcessID == _msg.ProcessID &&
			s.backendKeyData.SecretKey == _msg.SecretKey {

			s.currenthandlertag = HandlerTagCancelRequest
			err = s.h.HandlerCancelRequest(s.ctx_s, _msg)
		}
	}
	return err
}

func (s *server) isDebug(msg string, args ...any) {
	if s.log == nil && s.debug == false {
		return
	}
	s.log.Debug(msg, args...)
}
