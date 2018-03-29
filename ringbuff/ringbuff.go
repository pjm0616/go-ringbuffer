// Package ringbuff provides a minimalist ring buffer.
//
// Example:
//
//   buffer := New(10)
//   results := []int{}
//   expectedResults := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
//   buffer.Add(1)
//   buffer.Add(2)
//   buffer.Add(3)
//   buffer.Add(4)
//   buffer.Add(5)
//   buffer.Add(6)
//   buffer.Add(7)
//   buffer.Add(8)
//   buffer.Add(9)
//   buffer.Add(10)
//   buffer.Add(11)
//   buffer.ForEach(func(item interface{}) {
//      results = append(results, item.(int))
//   })
package ringbuff

type RingBuffer struct {
	items         []interface{}
	index         int
	highWaterMark int
}

// New creates a new RingBuffer capped at the specified size.
func New(size int) *RingBuffer {
	if size <= 0 {
		panic("invalid size for RingBuffer")
	}
	return &RingBuffer{
		items:         make([]interface{}, size),
		index:         0,
		highWaterMark: 0,
	}
}

// Add adds an item to the RingBuffer.
func (buffer *RingBuffer) Add(item interface{}) {
	buffer.items[buffer.index] = item
	// Update highWaterMark
	if buffer.index >= buffer.highWaterMark {
		buffer.highWaterMark = buffer.index + 1
	}
	// Advance index
	buffer.index++
	if buffer.index >= len(buffer.items) {
		buffer.index = 0
	}
}

// ForEach iterates over the RingBuffer starting with the oldest item.
func (buffer *RingBuffer) ForEach(fn func(interface{})) {
	if buffer.highWaterMark == 0 {
		// empty
		return
	}
	start := buffer.index - buffer.highWaterMark
	if start < 0 {
		// wrap around
		start += len(buffer.items)
	}
	for i := start; i < start+buffer.highWaterMark; i++ {
		index := i
		if index >= buffer.highWaterMark {
			// wrap around
			index = index - buffer.highWaterMark
		}
		fn(buffer.items[index])
	}
}

// Size reports how many items are currently held in the RingBuffer.
func (buffer *RingBuffer) Size() int {
	return buffer.highWaterMark
}

func (buffer *RingBuffer) Get(idx int) (interface{}, bool) {
	if idx < 0 || idx >= buffer.highWaterMark {
		return nil, false
	}
	ringIdx := buffer.index - buffer.highWaterMark
	if ringIdx < 0 {
		// wrap around
		ringIdx += len(buffer.items)
	}
	ringIdx += idx
	if ringIdx >= buffer.highWaterMark {
		// wrap around
		ringIdx = ringIdx - buffer.highWaterMark
	}
	return buffer.items[ringIdx], true
}
