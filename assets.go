package main

import "encoding/base64"

// TODO: the icons here are PNG, which work on linux/darwin, but probably not on
// windows as it wants ICO.

// systrayIconOff is base64-encoded from https://github.com/ipfs-shipyard/ipfs-desktop/blob/master/assets/icons/tray/off-22Template.png
const systrayIconOff = `
iVBORw0KGgoAAAANSUhEUgAAABIAAAAWCAYAAADNX8xBAAAAAXNSR0IArs4c6QAAAERlWElmTU0A
KgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAEqADAAQAAAAB
AAAAFgAAAABvKRSVAAABe0lEQVQ4Ec2U2U7DMBBFwybEviOBhFT+/6t4QkChQFnEzjlWJnLauEF9
4kpXduyZ6xmPJ1X137DQE9Ai+3u1zYjxu2Q/S2gTpyO4Ujt/MN7Ap/q7NXQJrWKhwDrU+RqKY6jo
C1TwDTbIhZZYPYC70BTuoOn8QKGtae5DU76Ht/ALpk3HLeiJij3CIfyEXVhm8RBuQ0WMeKyyUEiR
CLskoq17pqatPvpWqufwXs6hET3ASItpguntQCOKINKGikJVjS7hGjRsq/YOIzrXT6FCFiFsTW88
GZHhXsA49UwjKDzMIngnU9FOCmGT0rEiCpiCokJnU05VciFHl1Ds63AFfVfCeRGtCyta/WEjF3Ke
f/e5t+zD0cpYwQFM74JxFrQZQH1SVaP8zyxY6g3o5Vpqe8l7yi/b+zqBtokV9N5GsGkR58IIo598
V1ZPUfEK7UMfafRh81vRuAt2uX8AH2UOfyG2hw+yhZJQGNkyNrPwIfpg54aH9R04t/iU4y8EJkxf
Mb8umgAAAABJRU5ErkJggg==
`

func decodeAsset(src string) []byte {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		panic(err) // should not happen
	}
	return bs
}
