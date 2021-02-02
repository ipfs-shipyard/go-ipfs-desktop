package main

import (
	"os"
	"os/exec"
	"runtime"
	"time"
)

func browserCommands() [][]string {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	default: // Linux, etc
		cmds = append(cmds, []string{"xdg-open"})
	}
	// fallbacks
	cmds = append(cmds,
		[]string{"chromium"},
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"firefox"},
	)
	return cmds
}

func openBrowser(url string) bool {
	for _, args := range browserCommands() {
		cmd := exec.Command(args[0], append(args[1:], url)...)
		if cmd.Start() == nil && apparentSuccess(cmd) {
			return true
		}
	}
	return false
}

func apparentSuccess(cmd *exec.Cmd) bool {
	errc := make(chan error, 1)
	go func() { errc <- cmd.Wait() }()

	select {
	case <-time.After(2 * time.Second):
		return true
	case err := <-errc:
		return err == nil
	}
}
