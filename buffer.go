package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	if Debug == true {
		fmt.Printf("%+v\n", b)
	}
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
	if Debug == true {
		fmt.Printf("%+v\n", b)
	}
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
	files, err := ioutil.ReadDir(BufferTmpFolder)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, file := range files {
		bt, err := os.ReadFile(BufferTmpFolder + file.Name()) // just pass the file name
		if err != nil {
			log.Println(err)
			return err
		}
		var event Event
		err = json.Unmarshal(bt, &event)
		if err != nil {
			log.Println(err)
			return err
		}
		b.Add(&event)
	}
	return nil
}

func (b *eventBuffer) sendEvent() error {
	for {
		var events []*Event
		for {
			event, err := b.Get()
			if err != nil {
				break
			} else {
				events = append(events, event)
			}
		}
		if len(events) > 0 {
			jsonByte, err := json.Marshal(events)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return err
			}
			postUrl := "https://in-event.critsend.io/event/received/"
			r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonByte))
			if err != nil {
				log.Println(err)
				return err
			}
			client := &http.Client{}
			res, err := client.Do(r)
			if err != nil {
				log.Println(err)
				return err
			}

			defer res.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
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
