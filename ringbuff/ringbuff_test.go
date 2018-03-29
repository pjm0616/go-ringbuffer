package ringbuff

import (
	"testing"
)

func TestEmptyBuffer(t *testing.T) {
	buffer := New(10)
	results := []int{}
	expectedResults := []int{}
	buffer.ForEach(func(item interface{}) {
		results = append(results, item.(int))
	})
	if buffer.Size() != 0 {
		t.Errorf("Wrong size reported: %d", 0)
	}
	if !equals(results, expectedResults) {
		t.Errorf("Empty ringbuffer yielded wrong results: %v", results)
	}
}

func TestNonWrappedBuffer(t *testing.T) {
	buffer := New(10)
	results := []int{}
	expectedResults := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	buffer.Add(1)
	buffer.Add(2)
	buffer.Add(3)
	buffer.Add(4)
	buffer.Add(5)
	buffer.Add(6)
	buffer.Add(7)
	buffer.Add(8)
	buffer.Add(9)
	buffer.Add(10)
	buffer.ForEach(func(item interface{}) {
		results = append(results, item.(int))
	})
	if buffer.Size() != 10 {
		t.Errorf("Wrong size reported: %d", 0)
	}
	if !equals(results, expectedResults) {
		t.Errorf("Non-wrapped ringbuffer yielded wrong results: %v", results)
	}
}

func TestWrappedBuffer(t *testing.T) {
	buffer := New(10)
	results := []int{}
	expectedResults := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	buffer.Add(1)
	buffer.Add(2)
	buffer.Add(3)
	buffer.Add(4)
	buffer.Add(5)
	buffer.Add(6)
	buffer.Add(7)
	buffer.Add(8)
	buffer.Add(9)
	buffer.Add(10)
	buffer.Add(11)
	buffer.ForEach(func(item interface{}) {
		results = append(results, item.(int))
	})
	if buffer.Size() != 10 {
		t.Errorf("Wrong size reported: %d", 0)
	}
	if !equals(results, expectedResults) {
		t.Errorf("Wrapped ringbuffer yielded wrong results: %v", results)
	}
}

func equals(results []int, expectedResults []int) bool {
	if len(results) != len(expectedResults) {
		return false
	}
	for i, item := range expectedResults {
		if item != results[i] {
			return false
		}
	}
	return true
}

func TestWrappedBufferWithEvictHandler(t *testing.T) {
	buffer := New(10)

	// Set the handler.
	var evicted []int
	handler := func(item interface{}) {
		evicted = append(evicted, item.(int))
	}
	buffer.SetEvictHandler(handler)

	// Add items.
	results := []int{}
	expectedResults := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	buffer.Add(1)
	buffer.Add(2)
	buffer.Add(3)
	buffer.Add(4)
	buffer.Add(5)
	buffer.Add(6)
	buffer.Add(7)
	buffer.Add(8)
	buffer.Add(9)
	buffer.Add(10)
	buffer.Add(11)
	buffer.ForEach(func(item interface{}) {
		results = append(results, item.(int))
	})

	// Check the buffer contents.
	if buffer.Size() != 10 {
		t.Errorf("Wrong size reported: %d", 0)
	}
	if !equals(results, expectedResults) {
		t.Errorf("Wrapped ringbuffer yielded wrong results: %v", results)
	}
	// Check evicted items.
	if len(evicted) != 1 {
		t.Errorf("No item evicted")
	}
	if evicted[0] != 1 {
		t.Errorf("Wrong item has been evicted")
	}
}
