# simpleWebServiceTestGolang
Simple web service which allows users to store and retrieve plain text messages.

The service behaves as follows:
## post 
$ curl $domain/messages/ -d 'my test message to store'  

{"id":12345}	

## get
$ curl $domain/messages/12345  

my test message to store

## delete
curl -X "DELETE" $domain/messages/12345

removed

## put 
$ curl -X "PUT" $domain/messages/12345 -d 'my new test message to store'

{"id":12345}	

# Overview
Places the message (value) to store into a map with a unique id (key). The key is returned in the form {"id":12345}.	

When the a request for the value/message comes in, if it exists it will be retrieved, if it does not then there following error is sent. "message Id not found".

If a delete message comes in then the object is removed. 

# How to run:

Download the repo.
In terminal navigate to the route of downloaded folder
Type go install
Run go file from bin dir.

# How to Test:
There is a go test to ensure that update and retrieval from map works as expected.  
To run use the following command in the dowloaded directory.  

$ go test
  
From a mac simply paste these commands in to terminal. 

## post 
$ curl localhost:8080/messages/ -d 'my test message to store' 
 
{"id":12345}

## get
$ curl localhost:8080/messages/12345

my test message to store

## delete
$ curl -X "DELETE" localhost:8080/messages/12345

removed

## put 
$ curl -X "PUT" localhost:8080/messages/12345 -d 'my new test message to store'

{"id":12345}
