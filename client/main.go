package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

const (
	DISCONNECT_TIMEOUT = time.Second * 2
	COLOR_CURRENT_USER = "[#00ffff]"
	COLOR_OTHER_USERS = "[green]"
)

var (
	host = flag.String("host", "localhost", "server hostname")
	port = flag.String("port", "5003", "port on which server listens to")
	username = "user0"
)

func viewUpdater(ch chan string, msgView *tview.TextView, app *tview.Application) {
	for {
		msg := <-ch

		pos := strings.Index(msg, ":")
		sender := msg[:pos]
		text := msg[pos+1:]

		// if message from current user:
		if sender == username {
			msgView.Write([]byte(COLOR_CURRENT_USER))
		} else {
			msgView.Write([]byte(COLOR_OTHER_USERS))
		}
		msgView.Write([]byte(tview.Escape(fmt.Sprintf("[%s] ", sender))))
		msgView.Write([]byte("[white]"))
		msgView.Write([]byte(tview.Escape(text)))
		msgView.Write([]byte("\n"))
		app.Draw()
	}
}

func isEmptyString(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}

func allButColon(text string, ch rune) bool {
	return ch != ':'
}

func main() {
	flag.Parse()

	app := tview.NewApplication()

	// username view
	prompt := tview.NewInputField().
		SetLabel("USERNAME: ").
		SetAcceptanceFunc(allButColon)
	prompt.SetDoneFunc(func (key tcell.Key) {
			if text := prompt.GetText(); key == tcell.KeyEnter && !isEmptyString(text) {
				username = prompt.GetText()
				app.Stop()
			}
		})

	if err := app.SetRoot(prompt, true).Run(); err != nil {
		panic(err)
	}

	// connect
	conn := Connect(*host, *port)
	if conn == nil {
		os.Exit(1)
	}

	defer conn.Close()

	outChan := make(chan string)
	inChan := make(chan string)

	go MessageSender(conn, outChan)
	go MessageReceiver(conn, inChan)

	outChan <- username

	mainView := tview.NewTextView().
		SetDynamicColors(true)
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainView, 0, 1, false).
		AddItem(prompt, 1, 1, true)
	prompt.SetText("").SetLabel("> ").
		SetAcceptanceFunc(func (_ string, _ rune) bool { return true }).
		SetDoneFunc(func (key tcell.Key) {
			if text := prompt.GetText(); key == tcell.KeyEnter && !isEmptyString(text) {
				if text == "/quit" {
					app.Stop()
				} else {
					outChan <- text
					prompt.SetText("")
				}
			}
		})

	go viewUpdater(inChan, mainView, app)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	// disconnect
	Disconnect(conn)
}