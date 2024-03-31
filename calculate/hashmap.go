package calculate

import (
	"bytes"
	"sort"
	"sync"
)

type HashMap struct {
	Entries  []Station
	Capacity uint32
	Size     uint32
	mu       *sync.Mutex
}

func NewHashMap(capacity uint32) *HashMap {
	return &HashMap{
		Entries:  make([]Station, capacity),
		Capacity: capacity,
		Size:     0,
		mu:       &sync.Mutex{},
	}
}

func (h *HashMap) Set(key uint32, value Station) {
	h.mu.Lock()
	if h.Entries[key].Hash == key {
		h.Entries[key] = value
		h.mu.Unlock()
		return
	}

	h.Entries[key] = value
	h.Size++
	h.mu.Unlock()
}

func (h *HashMap) Get(key uint32) (Station, bool) {
	h.mu.Lock()
	val := h.Entries[key]
	h.mu.Unlock()
	if val.Count > 0 {
		return val, true
	}
	return Station{}, false
}

func (h *HashMap) Sort() {
	sort.Slice(h.Entries, func(i, j int) bool {
		return bytes.Compare(h.Entries[i].Name[:], h.Entries[j].Name[:]) < 0
	})
}

func (h *HashMap) Hash(bytes []byte) uint32 {
	var result uint32
	for _, b := range bytes {
		result = result*31 + uint32(b)
	}
	return result % h.Capacity
}
