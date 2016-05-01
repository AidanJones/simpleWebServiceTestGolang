package main

import (
	"fmt"
	"net/http"
	"strings"
)

var messageMap map[string]string
var umessageid int

func handler(w http.ResponseWriter, r *http.Request) {

	//TODO split handler in to 2 based on parsing of requies, might be advance MUX.

	var returnMessage string = ""

	//If there is some data sent in.
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
		var messageid string = addToMessageMap(message)
		fmt.Printf("Message ID " + messageid + "\n")
		returnMessage = "" + messageid
		//retun json  object with message id

		w.Header().Set("Content-Type", "application/json")

		//TODO create struct and parse this as json?
		fmt.Fprintf(w, "{\"id\":"+messageid+"}"+"\n")
	} else {
		//Now take the id from the path
		fmt.Printf("Got request.\n")
		fmt.Printf("Path " + r.URL.Path + "\n")
		var messageid string = strings.SplitAfter(r.URL.Path, "/")[2]
		fmt.Printf("Message ID " + messageid + "\n")
		var message string = retreiveFromMessageMap(messageid)
		if message != "" {
			returnMessage = "" + message
		} else {
			returnMessage = "message Id not found"
		}

		fmt.Fprintf(w, returnMessage+"\n")
	}

}

func addToMessageMap(value string) string {
	//TODO Could have uid in a concurency safe manner. Use a UID package?
	if messageMap == nil {
		messageMap = make(map[string]string)
	}
	umessageid = umessageid + 1
	var key string = fmt.Sprintf("%v", umessageid)
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
