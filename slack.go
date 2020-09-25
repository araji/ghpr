package main

import "fmt"

//Slacker account info to write to specific channel
type Slacker struct {
	username string
	password string
}

//SlackMessage json struct including channel and message to send
type SlackMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

/*SendMessage to slack channel
POST https://slack.com/api/chat.postMessage
Content-type: application/json
Authorization: Bearer xoxb-your-token
{
  "channel": "YOUR_CHANNEL_ID",
  "text": "Hello world :tada:"
}
*/
func (s *Slacker) SendMessage(slackMessage *SlackMessage) error {
	fmt.Printf("user %s is sending message : %s to channel %s \n", s.username, slackMessage.Text, slackMessage.Channel)
	return nil
}
