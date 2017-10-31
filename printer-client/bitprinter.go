package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"net"
)

import (
	"github.com/alexbrainman/printer"
	"github.com/lxn/walk"
	"github.com/zhouhui8915/go-socket.io-client"
)

func main() {
	appName := "KuponPrinter"
	appVersion := "v1.0.0"
	clientId := Hash(GetMacAddress())

	// SocketIO Client
	opts := &socketio_client.Options{}
	uri := "http://103.253.113.244:9292/"
	Client, error := socketio_client.NewClient(uri, opts)
	if error != nil {
		log.Fatal(error)
	}

	// Create MainWindow for attaching system tray icon
	mainWindow, error := walk.NewMainWindow()
	if error != nil {
		log.Fatal(error)
	}

	notifyIcon, error := walk.NewNotifyIcon()
	if error != nil {
		log.Fatal(error)
	}
	defer notifyIcon.Dispose() // Dispose notifyIcon after end of program execution

	icon, error := walk.NewIconFromResourceId(9)
	if error != nil {
		log.Fatal(error)
	}
	defer icon.Dispose() // Dispose icon after end of program execution

	if error := notifyIcon.SetIcon(icon); error != nil {
		log.Fatal(error)
	}

	if error := notifyIcon.SetToolTip(appName); error != nil {
		log.Fatal(error)
	}

	// notifyIcon actions

	// status action
	statusAction := walk.NewAction()
	if error := statusAction.SetText("Status Printer"); error != nil {
		log.Fatal(error)
	}
	statusAction.Triggered().Attach(func() {
		printerName, error := printer.Default()
		if error != nil {
			log.Fatal(error)
		}

		message := appName + " " + appVersion + "\n\n" + printerName + " terkoneksi ke server dengan id " + clientId
		walk.MsgBox(mainWindow, appName, message, walk.MsgBoxOK)
	})
	if error := notifyIcon.ContextMenu().Actions().Add(statusAction); error != nil {
		log.Fatal(error)
	}

	// exit action
	exitAction := walk.NewAction()
	if error := exitAction.SetText("Keluar"); error != nil {
		log.Fatal(error)
	}
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})
	if error := notifyIcon.ContextMenu().Actions().Add(exitAction); error != nil {
		log.Fatal(error)
	}

	// Display notifyIcon
	if error := notifyIcon.SetVisible(true); error != nil {
		log.Fatal(error)
	}

	// SocketIO client event handler
	Client.On("connection", func() {
		printerName, error := printer.Default()
		if error != nil {
			log.Fatal(error)
		}
		message := appName + " " + appVersion + "\n\n" + printerName + " terkoneksi ke server dengan id " + clientId

		if err := notifyIcon.ShowCustom(appName, message); err != nil {
			log.Fatal(err)
		}
	})

	Client.On("error", func() {
		log.Fatal("Error")
	})

	// Run loop
	mainWindow.Run()
}

func GetMacAddress() string {
	macAddress := ""
	ifs, _ := net.Interfaces()
	for _, v := range ifs {
		h := v.HardwareAddr.String()
		if len(h) == 0 {
			continue
		}
		macAddress += h
	}

	return macAddress
}

func Hash(raw string) string {
	hasher := sha1.New()
	hasher.Write([]byte(raw))
	byteArray := hasher.Sum(nil)
	return fmt.Sprintf("%x", byteArray)
}
