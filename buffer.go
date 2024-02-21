package event

import (
	"fmt"
	"sync"
)

type Buffer struct {
	data    []*Event
	maxSize int
	mutex   sync.Mutex
}

// Add adds an item to the buffer
func (b *Buffer) Add(item *Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if buffer is full
	if len(b.data) >= b.maxSize {
		fmt.Println("Buffer is full. Cannot add more items.")
		return
	}

	b.data = append(b.data, item)
}

// Get retrieves and removes an item from the buffer
func (b *Buffer) Get() *Event {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if buffer is empty
	if len(b.data) == 0 {
		fmt.Println("Buffer is empty. Returning empty event.")
		return &Event{}
	}

	// Get and remove the first item from the buffer
	item := b.data[0]
	b.data = b.data[1:]
	return item
}
