package commands

import "github.com/slack-games/slack-client"

// PingCommand ping back
func PingCommand() slack.ResponseMessage {

	return slack.ResponseMessage{
		Text: "You lucky found ping page",
	}
}
