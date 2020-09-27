package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

//JSONMessage holds message to send to slack
type JSONMessage struct {
	Text string `json:"text"`
}

//JSONColoredMessage colored message inner block
type JSONColoredMessage struct {
	Color     string   `json:"color"`
	MarkdwnIn []string `json:"mrkdwn_in"`
	Text      string   `json:"text"`
}

//ColoredMessage block to send to slack webhook
type ColoredMessage struct {
	Attachements []JSONColoredMessage `json:"attachments"`
}

//SendSlackMessage Using app wehbook
func SendSlackMessage(webhookURL, message, color string) {

	msgBlock := JSONColoredMessage{Color: color, MarkdwnIn: []string{"text", "fields"}, Text: message}
	slackMessage := &ColoredMessage{[]JSONColoredMessage{msgBlock}}
	//jsonMessage := JSONMessage{message}
	json, err := json.Marshal(slackMessage)
	log.Println(string(json))
	if err != nil {
		log.Printf("Error while serializing data :%s \n", err)
	}
	resp, err := http.Post(webhookURL, "Application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Printf("Error while submitting slack webhook request :%s \n", err)
	}
	defer resp.Body.Close()

}
