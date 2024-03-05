package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CritsendGo/modBuffer"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var Debug = false
var buffer *modBuffer.CSBuffer
var Token string
var ReadFolderInterval = 1 * time.Minute
var FolderEvent = "/tmp/event/"
var MaxEventInSend = 256
var IntervalEventSend = 250 * time.Millisecond

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
	buffer, err = modBuffer.NewBuffer(FolderEvent, 1024)
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
		nb := 0
		for {
			var event Event
			data, err := buffer.Get()
			if Debug {
				fmt.Println(data, err)
			}
			if err != nil || nb > MaxEventInSend {
				break
			}
			d, err := json.Marshal(data)
			err = json.Unmarshal(d, &event)
			if err != nil || nb > MaxEventInSend {
				break
			} else {
				events = append(events, &event)
			}
			nb++
		}
		if len(events) > 0 {
			jsonByte, err := json.Marshal(events)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			postUrl := "https://in-event.critsend.io/event/received/"
			if Debug {
				//fmt.Println(string(jsonByte), err)
			}
			r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonByte))
			if err != nil {
				log.Println(err)

			}
			client := &http.Client{}
			res, err := client.Do(r)
			resBody, err := ioutil.ReadAll(res.Body)
			if Debug {
				fmt.Println(string(resBody))
				fmt.Println(res)
			}
			defer res.Body.Close()
			if err != nil {
				log.Println(err)

			}
		}
		time.Sleep(IntervalEventSend)
	}

}
