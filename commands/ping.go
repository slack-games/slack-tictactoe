package commands

import "github.com/riston/slack-client"

// PingCommand ping back
func PingCommand() slack.ResponseMessage {
	return slack.ResponseMessage{
		Text:        "You lucky found ping page",
		Attachments: []slack.Attachment{},
	}
}
