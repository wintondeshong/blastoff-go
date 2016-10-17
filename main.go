package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"

	lib "github.com/wintondeshong/blastoff-go/lib"
)

func main() {

	ws, id := initialize()
	fmt.Println("blastoffbot ready, ^C exits")

	for {
		// read each incoming message
		m, err := lib.GetMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			// if so try to parse if
			parts := strings.Fields(m.Text)
			if len(parts) == 3 && parts[1] == "stock" {
				// looks good, get the quote and reply with the result
				go func(m lib.Message) {
					m.Text = lib.GetQuote(parts[2])
					lib.PostMessage(ws, m)
				}(m)
				// NOTE: the Message object is copied, this is intentional
			} else {
				m.Text = fmt.Sprintf("sorry, that does not compute\n")
				lib.PostMessage(ws, m)
			}
		}
	}
}

// Initialization
// --------------

func initialize() (*websocket.Conn, string) {
	if len(os.Args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: blastoff\n")
		os.Exit(1)
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file\n")
	}

	var slackToken = os.Getenv("SLACK_API_TOKEN")
	if len(slackToken) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable 'SLACK_API_TOKEN' is required\n")
		os.Exit(1)
	}

	// start a websocket-based Real Time API session
	return lib.SlackConnect(slackToken)
}
