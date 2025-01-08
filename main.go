package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
}

func main() {
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/send-message", handleSendMessage)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
