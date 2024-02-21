package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var eventBuffer *Buffer
var eventToken string

type Event struct {
	UserId     int
	Code       int
	CreateTime time.Time
	OriginTime time.Time
	Detail     string
	Id         string
	Recipient  string
	Source     string
}

func init() {
	eventBuffer = &Buffer{maxSize: 4096}
	eventBuffer.data = make([]*Event, 0)
}

// AddEvent Used to Add event
func AddEvent(e *Event) {
	eventBuffer.Add(e)
}

// SendEvent Used to Send event buffer if not empty
func SendEvent() error {
	var events []*Event
	for {
		newE := eventBuffer.Get()
		if newE.UserId == 0 {
			break
		} else {
			events = append(events, newE)
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
		r.Header = http.Header{"Authorization": []string{eventToken}}
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
	} else {
		//Nothing to send
	}
	return nil
}
