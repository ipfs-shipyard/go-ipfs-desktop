# Experiment: IPFS Desktop without Electron

This repo is an experiment created during a Hack Week (2021 Q1)
to explore if it is possible to provide most of IPFS Desktop functionality
without dependency on the Electron/Chromium Rube Goldberg machine.

Consider this experimental for now!


![go-ipfs-desktop_v0 0 1](https://user-images.githubusercontent.com/157609/107711238-96d89480-6cc7-11eb-862f-693ae02f8013.png)

## How to run this?

Pre-built artifacts can be found [here](https://github.com/ipfs-shipyard/go-ipfs-desktop/releases/latest).

They require `ipfs` (go-ipfs binary) to be present on `PATH` or in the same directory as `go-ipfs-desktop`.  
Set `IPFS_PATH` to override the default repo location.

This is a portable app, meaning no installation is required.

## Dependencies

### Linux

Make sure [libappindicator](https://launchpad.net/libappindicator) is present, eg. `apt-get install libappindicator3-dev`


## License

This project is dual-licensed under Apache 2.0 and MIT terms:

- Apache License, Version 2.0, ([LICENSE-APACHE](https://github.com/ipfs-shipyard/go-ipfs-desktop/blob/master/LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
- MIT license ([LICENSE-MIT](https://github.com/ipfs-shipyard/go-ipfs-desktop/blob/master/LICENSE-MIT) or http://opensource.org/licenses/MIT)
