package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/getlantern/systray"
)

func main() {
	log.SetPrefix("[go-ipfs-desktop] ")

	ipfsBinary, err := exec.LookPath("ipfs")
	if err != nil {
		log.Fatalf("could not locate the ipfs binary: %v", err)
	}
	log.Printf("using ipfs found at %s", ipfsBinary)

	daemon := trayDaemon{ipfsBinary: ipfsBinary}
	systray.Run(daemon.onReady, daemon.onExit)
}

type trayDaemon struct {
	// TODO: we might need some sort of lock, since clicks are async.
	ipfsBinary string
	ipfsCmd    *exec.Cmd
	webUIURL   string

	menuOpenWebUI *systray.MenuItem
	menuStartStop *systray.MenuItem
	menuVersion   *systray.MenuItem
	menuQuit      *systray.MenuItem
}

func (d *trayDaemon) onReady() {
	systray.SetIcon(systrayIconOff)
	systray.SetTooltip("IPFS Desktop")

	// This menu item is filled when we start the daemon below.
	d.menuStartStop = systray.AddMenuItem("Starting...", "Starting the IPFS daemon")

	systray.AddSeparator()

	d.menuOpenWebUI = systray.AddMenuItem("Open WebUI", "Open the WebUI in the default browser")
	d.menuOpenWebUI.Disable()

	d.menuVersion = systray.AddMenuItem("Version/bug TODO", "")
	d.menuQuit = systray.AddMenuItem("Quit", "Stop the IPFS daemon and quit the app")

	go func() {
		for {
			if err := d.handleClick(); err != nil {
				log.Printf("error: %v", err)
			}
		}
	}()

	// Start the daemon.
	d.startIPFS()
}

func (d *trayDaemon) onExit() {
	d.stopIPFS()
	log.Println("exiting")
}

func (d *trayDaemon) handleClick() error {
	select {
	case <-d.menuStartStop.ClickedCh:
		if d.ipfsRunning() {
			return d.stopIPFS()
		} else {
			return d.startIPFS()
		}

	case <-d.menuOpenWebUI.ClickedCh:
		return openBrowser(d.webUIURL)

	case <-d.menuVersion.ClickedCh:
		return fmt.Errorf("TODO")

	case <-d.menuQuit.ClickedCh:
		systray.Quit()
		return nil
	}
}

func (d *trayDaemon) ipfsRunning() bool {
	return d.ipfsCmd != nil
}

func (d *trayDaemon) startIPFS() error {
	log.Println("starting the IPFS daemon")

	d.menuStartStop.SetTitle("Starting...")
	d.menuStartStop.SetTooltip("Starting the IPFS daemon")
	d.menuStartStop.Disable() // while in progress

	cmd := exec.Command(d.ipfsBinary, "daemon")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// save stderr in case we crash
	// TODO: use stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("got line", line, len(line))
		webUIPrefix := "WebUI: "
		if strings.HasPrefix(line, webUIPrefix) {
			d.webUIURL = line[len(webUIPrefix):]
		} else if strings.Contains(line, "Daemon is ready") {
			// The daemon has started, we're done.
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	log.Println("IPFS daemon ready")

	d.ipfsCmd = cmd

	if d.webUIURL != "" {
		d.menuOpenWebUI.Enable()
	}

	d.menuStartStop.SetTitle("Stop IPFS")
	d.menuStartStop.SetTooltip("Stop the IPFS daemon")
	d.menuStartStop.Enable()
	systray.SetIcon(systrayIconOn)
	return nil
}

func (d *trayDaemon) stopIPFS() error {
	log.Println("stopping the IPFS daemon")

	d.menuStartStop.SetTitle("Stopping...")
	d.menuStartStop.SetTooltip("Stopping the IPFS daemon")
	d.menuStartStop.Disable() // while in progress
	d.menuOpenWebUI.Disable()

	// TODO: Send os.Interrupt first, giving ipfs a few seconds to
	// gracefully shut down. How to do that on Windows?
	d.ipfsCmd.Process.Kill()
	d.ipfsCmd.Wait()

	d.ipfsCmd = nil

	d.menuStartStop.SetTitle("Start IPFS")
	d.menuStartStop.SetTooltip("Start the IPFS daemon")
	d.menuStartStop.Enable()
	systray.SetIcon(systrayIconOff)
	return nil
}
