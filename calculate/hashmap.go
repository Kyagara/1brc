package calculate

import (
	"bytes"
	"sort"
)

type HashMap struct {
	Shards [][]Station
}

func NewHashMap(shards int) *HashMap {
	store := make([][]Station, shards)
	for i := range store {
		store[i] = make([]Station, hashMapSize)
	}

	return &HashMap{
		Shards: store,
	}
}

func (h *HashMap) Set(shard int, name []byte, temperature float32) {
	hash := h.Hash(name)

	// More '/' = more cpu time
	// Any first read or write to a field from Station is really slow,
	// I believe it has to do with cache misses

	station := &h.Shards[shard][hash] ///
	if station.Count == 0 {           ////////
		station.Hash = hash
		station.Name = name
	}

	station.Count++ //
	station.Total += temperature

	if temperature < station.Min {
		station.Min = temperature
	}

	if temperature > station.Max {
		station.Max = temperature
	}
}

func (h *HashMap) Sort() []Station {
	merged := make(map[uint32]Station)
	for _, shard := range h.Shards {
		for _, station := range shard {
			if station.Count > 0 {
				val, ok := merged[station.Hash]
				if !ok {
					merged[station.Hash] = station
					continue
				}

				val.Total += station.Total
				val.Count += station.Count
				if station.Min < val.Min {
					val.Min = station.Min
				}
				if station.Max > val.Max {
					val.Max = station.Max
				}
				merged[station.Hash] = val
			}
		}
	}

	sorted := make([]Station, 0, len(merged))
	for _, val := range merged {
		sorted = append(sorted, val)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return bytes.Compare(sorted[i].Name[:], sorted[j].Name[:]) < 0
	})

	return sorted
}

func (h *HashMap) Hash(bytes []byte) uint32 {
	var result uint32
	for _, b := range bytes {
		result = result*31 + uint32(b)
	}
	return result % hashMapSize
}
