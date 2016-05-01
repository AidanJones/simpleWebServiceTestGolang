package main

import "testing"

func TestAddAndRemoveFromMessageMap(t *testing.T) {

	message := "test message"
	messageid := addToMessageMap(message)

	if messageid == "" {
		t.Errorf("Error adding message")
	}

	retmessage := retreiveFromMessageMap(messageid)
	if message != retmessage {
		t.Errorf("Error retreiving message")
	}

}
