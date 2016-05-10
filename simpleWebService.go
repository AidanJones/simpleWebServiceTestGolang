package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/gorilla/mux"
)

var messageMap map[string]string // Map to store messages
var umessageid uint64            //counter for unique message id
var mu = &sync.Mutex{}           // Mutex used in lock of the messageMap, just used when updating the map.

type messageIdStruct struct {
	MessageId string `json:"id"`
}

func handlePostMessage(w http.ResponseWriter, r *http.Request) {

	//If there is some data sent in as a post message.
	//if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
	fmt.Printf("Got input new method.\n")

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

	mis := messageIdStruct{messageid}

	if err := json.NewEncoder(w).Encode(mis); err != nil {
		panic(err)
	}
}

func handleGetMessage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var messageid string
	messageid = vars["messageId"]

	//If there is a request for a message
	fmt.Printf("Got request.\n")
	fmt.Printf("Path " + r.URL.Path + "\n")
	//Take the id from the path
	//var messageid string = strings.SplitAfter(r.URL.Path, "/")[2]
	fmt.Printf("Message ID " + messageid + "\n")
	//Retreive message from map.
	message := retreiveFromMessageMap(messageid)
	if message == "" {
		//If no message to retreive send error message
		message = "message Id not found"
	}
	fmt.Printf("Message: " + message + "\n")
	fmt.Fprintf(w, message+"\n")

}

func addToMessageMap(message string) string {
	//Using mutex to lock the atomic increment of id and addition to map.
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}

	atomic.AddUint64(&umessageid, 1)
	var key string = fmt.Sprintf("%v", umessageid)
	mu.Unlock()

	messageMap[key] = message

	return key
}

func retreiveFromMessageMap(key string) string {
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}
	mu.Unlock()
	message := messageMap[key]
	return message
}

func main() {

	messageMap = make(map[string]string)
	umessageid = 12344

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/messages/{messageId}", handleGetMessage).Methods("GET")
	router.HandleFunc("/messages/", handlePostMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
