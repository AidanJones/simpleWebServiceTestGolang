package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"html/template"

	"github.com/gorilla/mux"
)

var messageMap map[string]string // Map to store messages
var umessageid uint64            //counter for unique message id
var mu = &sync.Mutex{}           // Mutex used in lock of the messageMap, just used when updating the map.

//TODO move this out to a file.
const allMessageTmpl = `
<body>

<div id="header">
<h1>All Messages</h1>
</div>
{{range .}}
	Message Id:"{{.MessageId}}"  Message:"{{.Message}}" <br>
{{end}}
</body>
`

type messageIdStruct struct {
	MessageId string `json:"id"`
	Message   string `json:"-"`
}

//If there is some data sent in as a post message.
func handlePostMessage(w http.ResponseWriter, r *http.Request) {
	//TODO do some sanity checks on the body or message passed in.
	//TODO use this as part of validation  if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
	fmt.Printf("Got input new method.\n")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	message := string(body)
	fmt.Printf("Message sent in " + message + "\n")
	//Add message to the map....
	var messageid string = addToMessageMap(message)
	fmt.Printf("Message ID " + messageid + "\n")

	//return json  object with message id

	mis := messageIdStruct{messageid, message}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(mis); err != nil {
		panic(err)
	}
}

func handleGetMessage(w http.ResponseWriter, r *http.Request) {
	//TODO validate the message id to make sure it meets the standards.
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

func handleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	//TODO validate the message id to make sure it meets the standards.
	//TODO can check message id exists before deleting to give more informative return value.
	vars := mux.Vars(r)
	var messageid string
	messageid = vars["messageId"]

	removeFromMessageMap(messageid)

	fmt.Fprintf(w, "removed"+"\n")

}

func handlePutMessage(w http.ResponseWriter, r *http.Request) {
	//TODO validate the message id to make sure it meets the standards.
	//TODO can check message id exists before deleting to give more informative return value.
	vars := mux.Vars(r)
	var messageid string
	messageid = vars["messageId"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	message := string(body)

	mis := messageIdStruct{messageid, message}
	putInMessageMap(mis)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(mis); err != nil {
		panic(err)
	}

}

func Index(w http.ResponseWriter, r *http.Request) {
	//TODO add some more html or link to the readme for the project....
	fmt.Fprintln(w, "Welcome!")
}

func putInMessageMap(mis messageIdStruct) {
	//Using mutex to lock the atomic increment of id and addition to map.
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}

	messageMap[mis.MessageId] = mis.Message

	mu.Unlock()

}

func removeFromMessageMap(key string) {
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}
	delete(messageMap, key)
	mu.Unlock()

}

func addToMessageMap(message string) string {
	//Using mutex to lock the atomic increment of id and addition to map.
	mu.Lock()
	if messageMap == nil {
		messageMap = make(map[string]string)
	}

	var key string
	//This ensures that the key is free, needed as put allows an id to be used out of sequence.
	var keyPresent bool = true
	for keyPresent {
		atomic.AddUint64(&umessageid, 1)
		key = fmt.Sprintf("%v", umessageid)
		message := messageMap[key]
		if message == "" {
			keyPresent = false
		}
	}

	messageMap[key] = message
	mu.Unlock()

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

func handleGetAllMessages(w http.ResponseWriter, r *http.Request) {
	//TODO add html buton to delete
	//TODO html button an pre filled text box to update message

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//mu.Lock()
	if len(messageMap) > 0 {
		allMessages := make([]messageIdStruct, 0)

		for k, v := range messageMap {
			//fmt.Println("k:", k, "v:", v)
			mis := messageIdStruct{k, v}
			allMessages = append(allMessages, mis)
		}
		//mu.Unlock()

		t := template.Must(template.New("allMessageTmpl").Parse(allMessageTmpl))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, allMessages)

		//		w.WriteHeader(http.StatusOK)
		//		if err := json.NewEncoder(w).Encode(allMessages); err != nil {
		//			panic(err)
		//		}
	} else {
		fmt.Fprintln(w, "No Messages!")
	}
}

func main() {
	//todo make it so that this can support arbitary message queue names?

	//todo replace map with db?
	messageMap = make(map[string]string)
	umessageid = 12344

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/messages/{messageId}", handleGetMessage).Methods("GET")
	router.HandleFunc("/messages/{messageId}", handleDeleteMessage).Methods("DELETE")
	router.HandleFunc("/messages/{messageId}", handlePutMessage).Methods("PUT")
	router.HandleFunc("/messages/", handlePostMessage).Methods("POST")
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/messages/", handleGetAllMessages).Methods("GET")

	//TODO handle for html request

	//TODO split some of the code in to packages
	log.Fatal(http.ListenAndServe(":8080", router))
}
