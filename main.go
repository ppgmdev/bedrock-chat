package main

import (
	"bedrock-chat/bedrock"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

var client *bedrockruntime.Client
var ctx context.Context

func init() {
    region := flag.String("region", "us-east-1", "AWS Region")
    ctx = context.Background()
    sdkconfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
    if err != nil {
        fmt.Println("couldn't load default sdk configuration")
        fmt.Println(err)
        return
    }
    client = bedrockruntime.NewFromConfig(sdkconfig)
}

type Message struct {
    Text string `json:"message"`
    User string `json:"user"`
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "not a POST request", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)

    if err != nil {
        http.Error(w, "error reading request body", http.StatusInternalServerError)
        return
    }

    defer r.Body.Close()

    var msg Message

    err = json.Unmarshal(body, &msg)
    if err != nil {
        http.Error(w, "error parsing JSON", http.StatusBadRequest)
        fmt.Println(string(body))
        fmt.Println(err)
        return
    }

    if msg.Text == "" {
        http.Error(w, "no message provided", http.StatusBadRequest)
        return
    }

    if msg.User == "" {
        http.Error(w, "no user provided", http.StatusBadRequest)
        return
    }

    fmt.Println("Request method:", r.Method)
    fmt.Println("Message:", msg.Text)
    
    bedrockConverse := bedrock.BedrockConverse{
        Message: msg.Text,
        Model: "amazon.nova-micro-v1:0",
    }

    response, err := bedrockConverse.NewMessage(ctx, client)
    if err != nil {
        fmt.Println("error calling bedrock")
        fmt.Println(err)
        http.Error(w, "Error calling bedrock", http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Bedrock: %s", response)
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    go http.Handle("/", fs)
    go http.HandleFunc("/send-message", handleSendMessage)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
