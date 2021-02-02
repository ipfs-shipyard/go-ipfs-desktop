package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"

	"github.com/getlantern/systray"
)

var ipfsBinary string

func main() {
	ipfsBinary, err := exec.LookPath("ipfs")
	if err != nil {
		log.Fatalf("could not locate the ipfs binary: %v", err)
	}
	log.Printf("using ipfs found at %s", ipfsBinary)

	daemon := trayDaemon{ipfsBinary: ipfsBinary}
	systray.Run(daemon.onReady, daemon.onExit)
}

type trayDaemon struct {
	ipfsBinary string

	menuDaemon  *systray.MenuItem
	menuVersion *systray.MenuItem
	menuQuit    *systray.MenuItem
}

func (d *trayDaemon) onReady() {
	icon, err := base64.StdEncoding.DecodeString(systrayIconOff)
	if err != nil {
		panic(err) // should not happen
	}

	systray.SetIcon(icon)
	systray.SetTooltip("IPFS Desktop")

	d.menuDaemon = systray.AddMenuItem("Start IPFS", "Start the IPFS daemon")

	systray.AddSeparator()

	d.menuVersion = systray.AddMenuItem("Version/bug TODO", "")
	d.menuQuit = systray.AddMenuItem("Quit", "Stop the IPFS daemon and quit the app")

	go func() {
		for {
			if err := d.handleClick(); err != nil {
				log.Printf("error: %v", err)
			}
		}
	}()
}

func (d *trayDaemon) handleClick() error {
	select {
	case <-d.menuDaemon.ClickedCh:
		return fmt.Errorf("TODO")

	case <-d.menuQuit.ClickedCh:
		systray.Quit()
		return nil

	case <-d.menuVersion.ClickedCh:
		return fmt.Errorf("TODO")
	}
}

func (d *trayDaemon) onExit() {
	log.Println("exiting")
}
