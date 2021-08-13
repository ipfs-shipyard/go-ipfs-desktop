package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	log.SetPrefix("[go-ipfs-desktop] ")

	ipfsBinary, err := locateIPFS()
	if err != nil {
		log.Fatalf("could not locate the ipfs binary: %v", err)
	}
	log.Printf("using ipfs found at %s", ipfsBinary)

	daemon := trayDaemon{
		daemonSemaphore: make(chan struct{}, 1),
		ipfsBinary:      ipfsBinary,
	}
	systray.Run(daemon.onReady, daemon.onExit)
}

// locateIPFS attempts to find the "ipfs" binary. It first checks the directory
// where go-ipfs-desktop is being run from, and then falls back to $PATH.
func locateIPFS() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", nil
	}
	execDir := filepath.Dir(execPath)
	for _, ext := range []string{"", ".exe"} {
		ipfsBinary := filepath.Join(execDir, "ipfs"+ext)
		if _, err := os.Stat(ipfsBinary); err == nil {
			return ipfsBinary, nil
		}
	}
	return exec.LookPath("ipfs")
}

type trayDaemon struct {
	// daemonSemaphore controls access to ipfsCmd and starting/stopping the
	// daemon, to make sure only one goroutine does that at a time.
	//
	// In general, use grabDaemon and releaseDaemon instead of the semaphore
	// directly.
	//
	// Even if we are careful to only call startIPFS and stopIPFS from a
	// single goroutine, this adds a layer of safety to ensure that we never
	// introduce memory races. Plus, we'll likely need more goroutines in
	// the future, such as those handling ^C or other user inputs.
	daemonSemaphore chan struct{}

	ipfsBinary string
	ipfsCmd    *exec.Cmd
	ourStop    bool // if true, we stopped the last daemon
	webUIURL   string

	menuStartStop *systray.MenuItem
	menuOpenWebUI *systray.MenuItem
	menuReportBug *systray.MenuItem
	menuQuit      *systray.MenuItem
}

func (d *trayDaemon) grabDaemon() bool {
	select {
	case d.daemonSemaphore <- struct{}{}:
		return true
	default:
		return false
	}
}

func (d *trayDaemon) releaseDaemon() {
	select {
	case <-d.daemonSemaphore:
	default:
		panic("called releaseDaemon without grabDaemon?")
	}
}

func (d *trayDaemon) onReady() {
	if !d.grabDaemon() {
		panic("we must grab the daemon semaphore when we start")
	}
	defer d.releaseDaemon()

	systray.SetIcon(systrayIconOff)
	systray.SetTooltip("IPFS Desktop")

	// This menu item is filled when we start the daemon below.
	d.menuStartStop = systray.AddMenuItem("Starting...", "Starting the IPFS daemon")

	systray.AddSeparator()

	d.menuOpenWebUI = systray.AddMenuItem("Open WebUI", "Open the WebUI in the default browser")
	d.menuOpenWebUI.Disable()

	d.menuReportBug = systray.AddMenuItem("Report a bug", "")
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
	if !d.grabDaemon() {
		// We are starting or stopping the daemon right now.
		// Directly wait for that to finish, up to a timeout.
		timeout := 5 * time.Second
		log.Printf("waiting for the daemon semaphore for up to %v...", timeout)

		select {
		case d.daemonSemaphore <- struct{}{}:
		case <-time.After(timeout):
			log.Printf("timed out; exiting")
			return
		}
	}
	// If the IPFS daemon is running, stop it.
	if d.ipfsCmd != nil {
		d.stopIPFS()
	}
	log.Printf("exiting")
}

func (d *trayDaemon) handleClick() error {
	select {
	case <-d.menuStartStop.ClickedCh:
		if !d.grabDaemon() {
			return fmt.Errorf("refusing to start/stop IPFS since it is already in progress")
		}
		defer d.releaseDaemon()
		if d.ipfsRunning() {
			return d.stopIPFS()
		} else {
			return d.startIPFS()
		}

	case <-d.menuOpenWebUI.ClickedCh:
		if !openBrowser(d.webUIURL) {
			return fmt.Errorf("could not open web UI with a browser")
		}
		return nil

	case <-d.menuReportBug.ClickedCh:
		// TODO: We could pre-fill this in the future, with e.g. the
		// version of this program as well as ipfs.
		if !openBrowser("https://github.com/ipfs-shipyard/go-ipfs-desktop/issues/new?template=bug_report.md") {
			return fmt.Errorf("could not open github with a browser")
		}
		return nil

	case <-d.menuQuit.ClickedCh:
		systray.Quit()
		return nil
	}
}

func (d *trayDaemon) ipfsRunning() bool {
	return d.ipfsCmd != nil
}

func (d *trayDaemon) startIPFS() error {
	if d.ipfsCmd != nil {
		panic("we tried to start the ipfs daemon when it's already running?")
	}
	d.ourStop = false // reset it

	log.Printf("starting the IPFS daemon")

	d.menuStartStop.SetTitle("Starting...")
	d.menuStartStop.SetTooltip("Starting the IPFS daemon")
	d.menuStartStop.Disable() // while in progress

	cmd := exec.Command(d.ipfsBinary, "daemon", "--migrate", "--routing", "dhtclient", "--enable-gc")
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
	ready := make(chan struct{}, 1)
	go func() {
		err := cmd.Wait()
		if !d.grabDaemon() {
			select {
			// The daemon started, and crashed at a later point.
			// The daemon is also being started or stopped.
			case <-ready:
				panic("daemon was stopped while grabbed?")
			// The daemon exited before it was ready.
			default:
				if err == nil {
					log.Printf("ipfs daemon exited without being ready, stderr below:\n%s", stderr.String())
				} else {
					log.Printf("ipfs daemon failed before being ready, stderr below: %v\n%s", err, stderr.String())
				}
			}
			return
		}
		defer d.releaseDaemon()

		if d.ourStop {
			// We triggered this daemon stop, so we don't need to do
			// anything else.

			// TODO: what if we send the daemon a ^C and it crashes
			// while shutting down?
			return
		}

		if err == nil {
			log.Printf("ipfs daemon exited unexpectedly, stderr below:\n%s", stderr.String())
		} else {
			log.Printf("ipfs daemon crashed unexpectedly, stderr below: %v\n%s", err, stderr.String())
		}

		d.menuStartStop.SetTitle("Start IPFS")
		d.menuStartStop.SetTooltip("Start the IPFS daemon")
		d.menuStartStop.Enable()
		systray.SetIcon(systrayIconOff)
		d.ipfsCmd = nil
	}()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		webUIPrefix := "WebUI: "
		if strings.HasPrefix(line, webUIPrefix) {
			d.webUIURL = strings.TrimSpace(line[len(webUIPrefix):])
		} else if strings.Contains(line, "Daemon is ready") {
			// The daemon has started, we're done.
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	log.Printf("IPFS daemon ready")

	if d.webUIURL == "" {
		log.Println("ipfs daemon started, but has no web UI?")
	} else {
		d.menuOpenWebUI.Enable()
	}

	d.menuStartStop.SetTitle("Stop IPFS")
	d.menuStartStop.SetTooltip("Stop the IPFS daemon")
	d.menuStartStop.Enable()
	systray.SetIcon(systrayIconOn)
	d.ipfsCmd = cmd
	ready <- struct{}{}
	return nil
}

func (d *trayDaemon) stopIPFS() error {
	if d.ipfsCmd == nil {
		panic("we tried to stop the ipfs daemon when it's not running?")
	}
	d.ourStop = true
	log.Printf("stopping the IPFS daemon")

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
