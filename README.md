# flydigictl

A utility for managing the configuration of [Flydigi](https://en.flydigi.com/) controllers.

Currently supported:

- Vader 3 Pro

Not tested but probably will work:

- Vader 3 Pro ONE PIECE
- Vader 3
- Vader 2
- Direwolf 2

Probably won't work:
- Vader 2
- Direwolf
- Others

## Installing
### Debian, Ubuntu and other Debian-based distros

You can download the artifact from the [latest actions run](https://github.com/pipe01/flydigictl/actions) and install it.

### Other distros

Install `libusb-1.0` and [Go](https://go.dev/) 1.21.5 or newer, then run `sudo make install` to install `flydigictl` and `flydigid`.

## Usage

This project consists of two parts: a daemon that runs in the background as root, and a command line utility that talks to this daemon through a DBus interface.
The daemon will be automatically started by SystemD when the DBus interface is requested.

To check if communication to the daemon works, run `flydigictl version`. You should see output like the following:

```
flydigictl version 0.0.1-8410f78
flydigid version 0.0.1-8410f78
```

To test communication with the controller, plug it in then run `flydigictl info`. You should see something like:

```
            Device : 80 (Vader 3 Pro)
Battery percentage : 0%
   Connection type : wireless
               CPU : wch (ch571)
```

Run `flydigictl help` to see what options the program has.

## Troubleshooting

### `flydigictl info` returns an error instead of information about the controller

- Make sure the controller is on DInput or XInput mode and not Bluetooth, Switch or others.
- Put the controller into XInput mode then check `lsusb`'s output. You should see a line that contains `ID 045e:028e Microsoft Corp. Xbox360 Controller`
- Unload the `xpad` module using `sudo modprobe -r xpad`. `flydigid` should do this automatically, but it may cause issues on some systems.

## Disclaimer

THIS PROGRAM IS PROVIDAD AS-IS AND MAKES NO EXPRESS OR IMPLIED WARRANTY OF ANY KIND. I AM NOT RESPONSIBLE FOR ANY POSSIBLE DAMAGE CAUSED TO THE DEVICE OR ANY ACCESSORIES.
