package event

import (
	"errors"
	"time"
)

var Debug = false
var buffer *eventBuffer
var Token string
var BufferTmpFolder string
var BufferErrFolder string
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
}

var (
	EventBoot      = 1
	EventVersion   = 2
	EventNbSuccess = 3
	EventNbError   = 4
)

func init() {
	buffer = &eventBuffer{maxSize: 4096}
	buffer.data = make([]*Event, 0)
	err := buffer.readEvent()
	if err != nil {
		// @TODO Show error on startup
	}
	// Init Async Reading Folder
	go asyncRead()
	go buffer.sendEvent()
}
func asyncRead() {
	for {
		err := buffer.readEvent()
		if err != nil {

		}
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

// WriteEvent Used to Write event IN DD
func WriteEvent(e *Event) error {

	return nil
}
