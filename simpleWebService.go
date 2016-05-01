package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
)

var messageMap map[string]string
var umessageid uint64
var mu = &sync.Mutex{}

func handler(w http.ResponseWriter, r *http.Request) {

	//TODO split handler in to two seperate methods based on parsing of URL, may be possible with advanced MUX.

	var returnMessage string = ""

	//If there is some data sent in as a post message.
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		fmt.Printf("Got input.\n")

		//It appears that the the data could also be taken from the body if the Parse form is not run.
		fmt.Printf("Request Form values.\n")
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("Error parsing form.\n")
		}
		var message string
		for k, v := range r.PostForm {
			fmt.Printf("  [%s]: \"%s\"\n", k, v)
			message = k
		}

		fmt.Printf("Message sent in " + message + "\n")
		//Add message to the map....
		var messageid string = addToMessageMap(message)
		fmt.Printf("Message ID " + messageid + "\n")

		//return json  object with message id
		//TODO create struct and parse this as json?
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"id\":"+messageid+"}"+"\n")
	} else {
		//If there is a reuest for a message
		fmt.Printf("Got request.\n")
		fmt.Printf("Path " + r.URL.Path + "\n")
		//Take the id from the path
		var messageid string = strings.SplitAfter(r.URL.Path, "/")[2]
		fmt.Printf("Message ID " + messageid + "\n")
		//Retreive message from map.
		var message string = retreiveFromMessageMap(messageid)
		if message != "" {
			//Send the stored message back
			returnMessage = "" + message
		} else {
			//Send error message
			returnMessage = "message Id not found"
		}

		fmt.Fprintf(w, returnMessage+"\n")
	}

}

func addToMessageMap(value string) string {
	//Using mutex to lock the atomic increment of id and addition to map.
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}

	atomic.AddUint64(&umessageid, 1)
	var key string = fmt.Sprintf("%v", umessageid)
	mu.Unlock()
	messageMap[key] = value
	return key
}

func retreiveFromMessageMap(key string) string {
	return messageMap[key]
}

func main() {

	messageMap = make(map[string]string)
	umessageid = 12344

	http.HandleFunc("/messages/", handler)
	http.ListenAndServe(":8080", nil)

}
