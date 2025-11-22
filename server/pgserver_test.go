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
