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
