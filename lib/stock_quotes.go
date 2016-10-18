package lib

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/websocket"
)

type StockHandler struct {

}


// Public Methods
// --------------

func (h StockHandler) CanHandleMessage(m Message) bool {
	parts := strings.Fields(m.Text)
	return len(parts) == 3 && parts[1] == "stock"
}

func (h StockHandler) HandleMessage(ws *websocket.Conn, m Message) {
	parts := strings.Fields(m.Text)

	m.Text = getQuote(parts[2])
	PostMessage(ws, m)
}


// Private Methods
// ---------------

// Get the quote via Yahoo. You should replace this method to something
// relevant to your team!
func getQuote(sym string) string {
	sym = strings.ToUpper(sym)
	url := fmt.Sprintf("http://download.finance.yahoo.com/d/quotes.csv?s=%s&f=nsl1op&e=.csv", sym)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	rows, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	if len(rows) >= 1 && len(rows[0]) == 5 {
		return fmt.Sprintf("%s (%s) is trading at $%s", rows[0][0], rows[0][1], rows[0][2])
	}
	return fmt.Sprintf("unknown response format (symbol was \"%s\")", sym)
}
