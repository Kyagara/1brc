package calculate

import (
	"bytes"
	"sort"
)

type HashMap struct {
	Entries  []Info
	Capacity uint32
	Size     uint32
}

func NewHashMap(capacity uint32) *HashMap {
	return &HashMap{Capacity: capacity, Size: 0, Entries: make([]Info, capacity)}
}

func (h *HashMap) Set(key uint32, value Info) {
	if h.Entries[key].Hash == key {
		h.Entries[key] = value
		return
	}

	h.Entries[key] = value
	h.Size++
}

func (h *HashMap) Get(key uint32) (Info, bool) {
	if h.Entries[key].Hash == key {
		return h.Entries[key], true
	}
	return Info{}, false
}

func (h *HashMap) Sort() {
	sort.Slice(h.Entries, func(i, j int) bool {
		return bytes.Compare(h.Entries[i].Name[:], h.Entries[j].Name[:]) < 0
	})
}

// Has conflitcs, few stations are missing
func (h *HashMap) Hash(bytes []byte) uint32 {
	var result uint32
	for _, b := range bytes {
		result = result*31 + uint32(b)
	}
	return result % h.Capacity
}
