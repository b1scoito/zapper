package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type RedirectMessages struct {
	WMClient        *whatsmeow.Client
	eventHandlerID  uint32
	groupsToSend    []string
	groupsToReceive []string
}

func (cli *RedirectMessages) registerEvents() {
	cli.eventHandlerID = cli.WMClient.AddEventHandler(cli.eventHandler)
}

func (cli *RedirectMessages) unregisterEvents() {
	cli.WMClient.RemoveEventHandler(cli.eventHandlerID)
}

func (cli *RedirectMessages) sendMessages(v *events.Message) {
	for _, group := range cli.groupsToSend {
		// Check if group is not the same as the one that sent the message
		if v.Info.Chat.User != group {
			go cli.WMClient.SendMessage(context.Background(), types.NewJID(group, types.GroupServer), v.Message)
		}
	}
}

func (cli *RedirectMessages) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		for _, group := range cli.groupsToReceive {
			if v.Info.Chat.User == group {
				cli.sendMessages(v)
			}
		}
	}
}

func main() {
	var groupsSend string
	var groupsReceive string

	// 120363122182986428
	flag.StringVar(&groupsSend, "sends", "", "List of groups to send messages to, separated by a comma. (e.g. 120363122182986428,120363120553285556)")

	// 120363120553285556
	flag.StringVar(&groupsReceive, "receives", "", "List of groups to receive messages from, separated by a comma. (e.g. 120363122182986428,120363120553285556)")
	flag.Parse()

	// If argument is null, exit
	if len(groupsSend) <= 0 || len(groupsReceive) <= 0 {
		fmt.Printf("Usage: %s -sends 120363122182986428 -receives 120363120553285556\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	dbLog := waLog.Stdout("Database", "INFO", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:wadata.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	cli := &RedirectMessages{WMClient: client}
	cli.registerEvents()
	cli.groupsToSend = strings.Split(groupsSend, ",")
	cli.groupsToReceive = strings.Split(groupsReceive, ",")

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	cli.unregisterEvents()
	client.Disconnect()
}
