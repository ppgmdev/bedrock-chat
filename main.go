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
    Model string `json:"model"`
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
    
    if msg.Model == "" {
        http.Error(w, "no user provided", http.StatusBadRequest)
        return
    }

    fmt.Println("Request method:", r.Method)
    fmt.Println("Message:", msg.Text)
    
    bedrockConverse := bedrock.BedrockConverse{
        Message: msg.Text,
    }
    
    switch msg.Model {
    case "nova-micro":
        bedrockConverse.Model = "amazon.nova-micro-v1:0"
    case "nova-lite":
        bedrockConverse.Model = "amazon.nova-lite-v1:0"
    case "nova-pro":
        bedrockConverse.Model = "amazon.nova-pro-v1:0"
    case "haiku-3":
        bedrockConverse.Model = "anthropic.claude-3-haiku-20240307-v1:0"
    case "claude-35-v2":
        bedrockConverse.Model = "us.anthropic.claude-3-5-sonnet-20241022-v2:0"
    case "haiku-35":
        bedrockConverse.Model = "us.anthropic.claude-3-5-haiku-20241022-v1:0"
    case "llama-33-70b":
        bedrockConverse.Model = "us.meta.llama3-3-70b-instruct-v1:0"
    case "mistral-small":
        bedrockConverse.Model = "mistral.mistral-small-2402-v1:0"
    case "mistral-7b":
        bedrockConverse.Model = "mistral.mistral-7b-instruct-v0:2"
    case "jamba-large":
        bedrockConverse.Model = "ai21.jamba-1-5-large-v1:0"
    case "jamba-mini":
        bedrockConverse.Model = "ai21.jamba-1-5-mini-v1:0"
    }


    response, err := bedrockConverse.NewMessage(ctx, client)
    if err != nil {
        fmt.Println("error calling bedrock")
        fmt.Println(err)
        http.Error(w, "Error calling bedrock", http.StatusInternalServerError)
        return
    }

    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)

}

func main() {
    fmt.Println("server listening on port 8080")
    fs := http.FileServer(http.Dir("./static"))
    go http.Handle("/", fs)
    go http.HandleFunc("/send-message", handleSendMessage)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
