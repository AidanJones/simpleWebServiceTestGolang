# simpleWebServiceTestGolang
Simple web service which allows users to store and retrieve plain text messages.


Develop a simple web service which allows users to store and retrieve plain text messages.

The service should behave as follows:

$ curl $domain/messages/ -d 'my test message to store' 
{"id":12345}	

$ curl $domain/messages/12345
my test message to store

# Overview
Places the message (value) to store into a map with a unique id (key). The key is returned in the form {"id":12345}.	

When the a request for the value comes in, if it exists it will be retrieved, if it does not then there following error is sent. "message Id not found".


# How to run:

Download the repo.
In terminal navigate to the route of downloaded folder
Type go install
Run go file from bin dir.

# How to Test:
From a mac simple paste these commands in to terminal. 

$ curl localhost:8080/messages/ -d 'my test message to store' 
{"id":12345}

$ curl localhost:8080/messages/12345