package event

import (
	"fmt"
	modBuffer "github.com/CritsendGo/modBuffer"
	"testing"
	"time"
)

var folderTest = "/tmp/eventSend/"
var size = 2
var Log = false
var BufferFolder = "/tmp/"

func TestAll(t *testing.T) {
	Debug = true
	Log = true
	var e1 = Event{
		UserId:     1,
		Code:       150,
		CreateTime: time.Now(),
		OriginTime: time.Now(),
		Detail:     "Test 1",
		Id:         "ddd",
		Recipient:  "none@none.com",
	}
	err := buffer.Add(e1)
	err = buffer.Add(e1)
	err = buffer.Add(e1)
	err = buffer.Add(e1)
	err = buffer.Add(e1)
	err = buffer.Add(e1)
	err = buffer.Add(e1)

	fmt.Println(err)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	time.Sleep(1 * time.Minute)
}

func addContent(bu *modBuffer.CSBuffer, val string) string {
	obj := Event{}
	err := bu.Add(obj)
	if Log == true {
		fmt.Println("ADDING CONTENT")

	}
	fmt.Println(err)
	return obj.Recipient
}
