package commands

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/slack-games/slack-client"
	"github.com/slack-games/slack-tictactoe"
	tttdatastore "github.com/slack-games/slack-tictactoe/datastore"
)

// StartCommand is command to start
func StartCommand(db *sqlx.DB, userID string) slack.ResponseMessage {
	var attachment slack.Attachment
	baseURL := os.Getenv("BASE_PATH")
	message := "There's already existing a game, you have to finish it before starting a new"

	// Try to get user last state
	state, err := tttdatastore.GetUserLastState(db, userID)

	if err != nil {
		// No state found
		if err == sql.ErrNoRows {
			now := time.Now().Unix()

			state = tttdatastore.State{
				State:        "000010000",
				TurnID:       userID,
				Mode:         "Start",
				FirstUserID:  "U000000000",
				SecondUserID: userID,
				ParentID:     "00000000-0000-0000-0000-000000000000",
				Created:      time.Now(),
			}

			if (now % 2) == 0 {
				state.FirstUserID = userID
				state.SecondUserID = "U000000000"
			}

			log.Println("Create a new state")
			stateID, err := tttdatastore.NewState(db, state)
			if err != nil {
				log.Fatalln("Could not create a new state", err)
			}

			message = "Created a new clean game state"
			attachment = slack.Attachment{
				Title:    "Last game state",
				Text:     "",
				Fallback: "Text fallback if image fails",
				ImageURL: fmt.Sprintf("%s/game/tictactoe/image/%s", baseURL, stateID),
				Color:    "#764FA5",
			}

			log.Println("New state id", stateID)
		} else {
			log.Println("Error could not get the user state")
		}
	} else if isGameOver(state) {
		state = tttdatastore.State{
			State:        "000000000",
			TurnID:       userID,
			Mode:         "Start",
			FirstUserID:  "U000000000",
			SecondUserID: userID,
			ParentID:     "00000000-0000-0000-0000-000000000000",
			Created:      time.Now(),
		}

		log.Println("Create a new state")
		stateID, err := tttdatastore.NewState(db, state)
		if err != nil {
			log.Fatalln("Could not create a new state", err)
		}

		message = "Created a new game state, last one is over"
		attachment = slack.Attachment{
			Title:    "New game state",
			Fallback: "Text fallback if image fails",
			ImageURL: fmt.Sprintf("%s/game/tictactoe/image/%s", baseURL, stateID),
			Color:    "#764FA5",
		}
	} else {
		attachment = slack.Attachment{
			Title:    "Last game state",
			Fallback: "Text fallback if image fails",
			ImageURL: fmt.Sprintf("%s/game/tictactoe/image/%s", baseURL, state.StateID),
			Color:    "#764FA5",
		}
	}

	return slack.ResponseMessage{
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
}

func isGameOver(state tttdatastore.State) bool {
	return state.Mode == fmt.Sprintf("%s", tictactoe.GameOverState) ||
		state.Mode == fmt.Sprintf("%s", tictactoe.WinState) ||
		state.Mode == fmt.Sprintf("%s", tictactoe.DrawState)
}
