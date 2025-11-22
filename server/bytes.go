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
