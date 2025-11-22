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
	"errors"
	"net"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/stretchr/testify/assert"
)

type Send func(pgproto3.FrontendMessage) error
type Receive func() (pgproto3.BackendMessage, error)

type _IO struct {
	Send    Send
	Receive Receive
	Write   func([]byte) (int, error)
	Read    func([]byte) (int, error)
}

type Func func(*_IO) error

func StartClient(t *testing.T, conn net.Conn, cmds map[string]Func) {
	var (
		io _IO
	)
	defer conn.Close()

	front := pgproto3.NewFrontend(conn, conn)
	defer func() {
		err := front.Flush()
		if err != nil {
			panic(err)
		}

	}()

	io.Send = func(fmsg pgproto3.FrontendMessage) error {
		if fmsg == nil {
			return errors.New("parameter fmsg can not be nil")
		}
		front.Send(fmsg)
		err := front.Flush()
		if err != nil {
			return err
		}
		return nil
	}
	io.Receive = func() (pgproto3.BackendMessage, error) {
		return front.Receive()
	}
	io.Write = conn.Write
	io.Read = conn.Read

	for _, cmd := range cmds {

		time.Sleep(2 * time.Second)

		err := cmd(&io)
		assert.Nil(t, err)
		time.Sleep(2 * time.Millisecond)
	}
}
