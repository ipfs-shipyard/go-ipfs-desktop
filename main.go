package main

import (
	"encoding/base64"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	icon, err := base64.StdEncoding.DecodeString(systrayIconOff)
	if err != nil {
		panic(err)
	}
	systray.SetIcon(icon)
	systray.SetTooltip("IPFS Desktop")
	mQuit := systray.AddMenuItem("Quit", "Stop the IPFS daemon and quit the app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// clean up here
}
