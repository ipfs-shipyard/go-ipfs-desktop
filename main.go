package main

import (
	"encoding/base64"
	"log"
	"os/exec"

	"github.com/getlantern/systray"
)

var ipfsBinary string

func main() {
	var err error
	ipfsBinary, err = exec.LookPath("ipfs")
	if err != nil {
		log.Fatalf("could not locate the ipfs binary: %v", err)
	}

	systray.Run(trayReady, trayExit)
}

func trayReady() {
	icon, err := base64.StdEncoding.DecodeString(systrayIconOff)
	if err != nil {
		panic(err)
	}

	systray.SetIcon(icon)
	systray.SetTooltip("IPFS Desktop")

	menuDaemon := systray.AddMenuItem("Start IPFS", "Start the IPFS daemon")

	systray.AddSeparator()

	menuVersion := systray.AddMenuItem("Version/bug TODO", "")
	menuQuit := systray.AddMenuItem("Quit", "Stop the IPFS daemon and quit the app")

	go func() {
		for {
			select {
			case <-menuDaemon.ClickedCh:
				println("TODO")

			case <-menuQuit.ClickedCh:
				systray.Quit()

			case <-menuVersion.ClickedCh:
				println("TODO")
			}
		}
	}()
}

func trayExit() {
}
