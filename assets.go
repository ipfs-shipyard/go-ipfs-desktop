package main

import "encoding/base64"

// TODO: use go:embed once Go 1.16 ships, then require that Go version to build

func decodeAsset(src string) []byte {
	bs, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		panic(err) // should not happen
	}
	return bs
}
