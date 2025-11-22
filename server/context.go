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
