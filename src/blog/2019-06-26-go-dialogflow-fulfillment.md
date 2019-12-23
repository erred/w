--- title
go dialogflow fulfillment
--- description
Writing DialogFlow fulfillment servers in Go
--- main


[DialogFlow](https://dialogflow.com/),
unified backend chatbot interface powered by ai

### Fulfillment

Calling your backend servers based on front chatbot input

### Go

First thing you notice,
there is only an SDK for client side Go,
**(sending user text to DialogFlow and getting the response)**,
and nothing for server side fulfillment.

No worries,
everything is standardized

### Webhook Request Response

The godoc for the message types can be found [here](https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2),

But the webhook is called with json,
and that it proto

It includes json tags,
but the json we're getting is encoded with the json parameter in the protobuf tag,
[jsonpb](https://godoc.org/github.com/golang/protobuf/jsonpb) to the rescue

### minimal code

```go
package main

import (
    "net/http"

    "github.com/golang/protobuf/jsonpb"
    pb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
        req := &pb.WebhookRequest{}
        err := jsonpb.Unmarshal(r.Body, req)
        if err != nil {
            // handle error
        }
        // Do something with req

        res := &pb.&pb.WebhookResponse{FulfillmentText:"I am fulfilled!"}

        m := &jsonpb.Marshaler{}
        err = m.Marshal(w, res)
        if err != nil {
            // handle error
        }
    })
}

```
