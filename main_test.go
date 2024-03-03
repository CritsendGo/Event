package event

import (
	"fmt"
	"github.com/CritsendGo/modBuffer"
	"os"
	"testing"
)

var folderTest = "/tmp/eventSend/"
var size = 2
var Log = false

func TestAll(t *testing.T) {
	Debug = false
	Log = true
	err := os.RemoveAll(folderTest)
	if err != nil {
		t.Fatal("Unable to clean folder buffer", err)
	}
	bu, err := modBuffer.NewBuffer(folderTest, size)
	if err != nil {
		t.Fatal("Unable to create buffer", err)
	}
	rep := addContent(bu, "First")
	if rep != "First" {
		t.Fatal("Add One Entry , BUFFER = ", 1, "FOLDER=", bu.SizeNew(), "ERROR:", err, "CONTENT", rep)
	}
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
