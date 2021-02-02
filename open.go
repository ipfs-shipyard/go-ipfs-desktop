package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Copied from https://gist.github.com/nanmu42/4fbaf26c771da58095fa7a9f14f23d27
func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return
}
