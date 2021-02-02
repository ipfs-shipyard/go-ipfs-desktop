## Native IPFS Desktop

This is an alternative to ipfs-desktop, showing a system tray icon to run an
IPFS daemon. There are several advantages to this implementation:

* No external dependencies like Electron
* Written in Go; can be bundled with go-ipfs in the future
* Easy to distribute as a mostly-static binary

Note that, by design, this project is not pure Go. It needs to make use of Cgo
to interact with the system.
