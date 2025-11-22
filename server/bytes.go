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
	"sync"

	maps "github.com/zolstein/sync-map"
)

var (
	mu sync.RWMutex

	_pools   maps.Map[uint16, *sync.Pool]
	_initial uint16 = 1

	BytesPure uint16 = 0
)

type Bytes struct {
	Cap   uint16
	Index uint16
	Data  []byte
}

func (b *Bytes) Get() {
	if b.Index != BytesPure {
		pool, ok := _pools.Load(b.Index)
		if !ok {
			panic("exbytes: Bytes.Index does not exits")
		}
		b.Data = pool.Get().([]byte)
	}
	if b.Data == nil {
		b.Data = make([]byte, 0, b.Cap)
	}
}

func (b *Bytes) Put() {
	if b.Data == nil {
		return
	}
	if b.Index == BytesPure {
		b.Data = nil
		return
	}

	pool, ok := _pools.Load(b.Index)
	if ok {
		clear(b.Data)
		pool.Put(b.Data)
	}
}

func (b *Bytes) Init() {
	if b.Index != BytesPure {
		if b.Cap == 0 {
			panic("Bytes.Cap cannot be zero")
		}
		_pools.LoadOrStore(b.Index, &sync.Pool{
			New: func() any {
				return make([]byte, 0, b.Cap)
			},
		})
		return
	}
	if len(b.Data) > 0 {
		return
	}
	b.Data = make([]byte, 0, b.Cap)
}

func Count() uint16 {
	i := _initial + 1
	return i
}
