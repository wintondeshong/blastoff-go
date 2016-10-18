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

	handlers := make([]lib.ISlackHandler, 1)
	handlers[0] = &lib.StockHandler{}

	for {
		m, err := lib.GetMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		if m.Type != "message" || ! strings.HasPrefix(m.Text, "<@" + id + ">") {
			continue
		}

		//var handler lib.ISlackHandler = lib.StockHandler{}
		handler := findHandler(handlers, m)
		if handler == nil {
			m.Text = fmt.Sprintf("Command not recognized\n")
			lib.PostMessage(ws, m)
			continue
		}

		go func(m lib.Message) {
			handler.HandleMessage(ws, m)
		}(m)
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

// Private Methods
// ---------------

func findHandler(handlers []lib.ISlackHandler, m lib.Message) lib.ISlackHandler {
	for i := 0; i < len(handlers); i++ {
		handler := handlers[i]

		if handler.CanHandleMessage(m) {
			return handler
		}
	}

	return nil
}
