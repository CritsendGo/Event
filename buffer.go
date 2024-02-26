package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type eventBuffer struct {
	data    []*Event
	maxSize int
	mutex   sync.Mutex
}

// SaveEvent Used to Save event in DD on JSON FORMAT
func (b *Event) saveEvent() (string, error) {
	fileName := fmt.Sprintln(time.Now().UnixMicro())
	filePath := BufferTmpFolder + fileName
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		log.Println("CREATE EVENT DD", err)
		return fileName, err
	}
	eBit, err := json.Marshal(b)
	if err != nil {
		log.Println("JSON EVENT DD", err)
		return fileName, err
	}
	l, err := f.Write(eBit)
	if err != nil {
		log.Println("WRITE EVENT DD", err)
		return fileName, err
	}
	fmt.Println(l, "bytes written successfully")
	return fileName, nil
}

// Add adds an item to the buffer
func (b *eventBuffer) Add(item *Event) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if buffer is full
	if len(b.data) >= b.maxSize {
		// Write on disk event return on success critical on error
		fileName, err := item.saveEvent()
		if err != nil {
			log.Println("BUFFER FULL AND UNABLE TO SAVE EVENT ON DISK", fileName)
			return
		}
		return
	}

	b.data = append(b.data, item)
}

// Get retrieves and removes an item from the buffer
func (b *eventBuffer) Get() (*Event, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	// Check if buffer is empty
	if len(b.data) == 0 {
		return &Event{}, errors.New("no more")
	}
	// Get and remove the first item from the buffer
	item := b.data[0]
	b.data = b.data[1:]
	return item, nil
}

// ReadEvent Used to Read event in DD
func (b *eventBuffer) readEvent() error {
	files, err := ioutil.ReadDir("/tmp/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// @TODO Unmarshal file content and add it to buffer

		// @TODO REMOVE
		fmt.Println(file.Name(), file.IsDir())
	}

	return nil
}

// ShutDown is used to save on disk the all buffer memory before shutdown and loose buffer
func (b *eventBuffer) ShutDown() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for {
		item, err := b.Get()
		// Check if buffer is empty
		if err != nil && item.Code == 0 {
			// No More Event
			break
		}
		// Not empty write to disk and do next
		fileName, errS := item.saveEvent()
		if errS != nil {
			log.Println("SHUTDOWN", "UNABLE TO SAVE EVENT ON DISK", fileName, errS)
		}
	}
}
