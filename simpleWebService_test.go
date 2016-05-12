package main

import "testing"

//TODO add more tests. Do the put out of order and ensure that duplicate id is not used.

//Test adding a message to the map and retreiving it.
func TestAddAndRetreiveFromMessageMap(t *testing.T) {

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

//test retreiving a message that is not in map
func TestRetreiveFromMessageMap(t *testing.T) {

	message := retreiveFromMessageMap("12345")
	if message != "" {
		t.Errorf("Unexpected return value")
	}

}
