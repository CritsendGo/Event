package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CritsendGo/modBuffer"
	"log"
	"net/http"
	"time"
)

var Debug = false
var buffer *modBuffer.CSBuffer
var Token string
var ReadFolderInterval = 1 * time.Minute

type Event struct {
	UserId     int
	Code       int
	CreateTime time.Time
	OriginTime time.Time
	Detail     string
	Id         string
	Recipient  string
	Source     string
	Tags       []string
}

var (
	EventBoot      = 1
	EventVersion   = 2
	EventNbSuccess = 3
	EventNbError   = 4
)

func init() {
	var err error
	modBuffer.Debug = Debug
	buffer, err = modBuffer.NewBuffer("", 1024)
	if err != nil {
		log.Println(err)
	}

	// Init Async Reading Folder
	go asyncRead()
	go sendEvent()
}
func asyncRead() {
	for {
		buffer.ScanFolder()
		time.Sleep(ReadFolderInterval)
	}
}

// AddEvent Used to Add event
func AddEvent(e *Event) error {
	if Token == "" {
		return errors.New("please set the Token var before calling AddEvent()")
	}
	buffer.Add(e)
	return nil
}
func sendEvent() {
	for {
		var events []*Event
		for {
			item, err := buffer.Get()
			if err != nil {
				break
			} else {
				event := item.(Event)
				events = append(events, &event)
			}
		}
		if len(events) > 0 {
			jsonByte, err := json.Marshal(events)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			postUrl := "https://in-event.critsend.io/event/received/"
			r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonByte))
			if err != nil {
				log.Println(err)

			}
			client := &http.Client{}
			res, err := client.Do(r)
			defer res.Body.Close()
			if err != nil {
				log.Println(err)

			}
		}

	}
	time.Sleep(100 * time.Millisecond)
}
