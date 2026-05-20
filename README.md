# flydigictl

`flydigictl` configures Flydigi controllers on Linux. It is split into two programs:

- `flydigid`: a root daemon that talks to the controller over USB.
- `flydigictl`: a user-facing CLI that talks to `flydigid` over system DBus.

## Supported Controllers

Tested or actively supported:

- Vader 3 Pro, device ID `80`
- Vader 4 Pro, device ID `105`

Known device IDs:

- `80`: Vader 3 Pro
- `81`: Vader 3 Pro ONE PIECE
- `85`: Vader 4
- `105`: Vader 4 Pro

Other Flydigi devices may work if they use the same USB protocol, but they are not verified.

## Requirements

- Linux with systemd and system DBus
- `libusb-1.0`
- Go 1.21.5 or newer
- A controller connected in XInput or DInput mode

Bluetooth, Switch mode, and other non-XInput/non-DInput modes are not supported by this tool.

## Build

Build both binaries from the repository root:

```bash
make bin-daemon bin-ctl
```

If Go is installed through a version manager such as `mise`, make sure that Go is on `PATH` first:

```bash
export PATH="/path/to/go/bin:$PATH"
make bin-daemon bin-ctl
```

## Install

The upstream `make install` target installs the DBus policy to `/etc/dbus-1/system.d`. On Arch Linux, install it under `/usr/share/dbus-1/system.d` instead.

Arch Linux manual install:

```bash
sudo cp flydigictl flydigid /usr/bin/
sudo cp etc/flydigid.service /usr/lib/systemd/system/flydigid.service
sudo cp etc/flydigid.conf /usr/share/dbus-1/system.d/flydigid.conf
sudo cp etc/com.pipe01.flydigi.Gamepad.service /usr/share/dbus-1/system-services/com.pipe01.flydigi.Gamepad.service
sudo systemctl daemon-reload
```

Other distributions may work with:

```bash
sudo PATH="$PATH" make install
```

## Service

`flydigid` is DBus-activatable. In normal use, running `flydigictl` starts the daemon automatically.

Useful service commands:

```bash
sudo systemctl start flydigid.service
sudo systemctl stop flydigid.service
sudo systemctl restart flydigid.service
sudo journalctl -u flydigid -f
```

## Controller Protocol Mode

`flydigid` defaults to `auto` mode, which tries DInput/HID first and falls back to XInput. You can force a mode with either a daemon flag or an environment variable:

```bash
flydigid --mode auto
flydigid --mode dinput
flydigid --mode xinput
FLYDIGI_MODE=dinput flydigid
```

Accepted mode names are `auto`, `dinput`, and `xinput`. `hid` and `directinput` are accepted aliases for `dinput`.

For the systemd service, run `sudo systemctl edit flydigid.service` and add a drop-in override:

```ini
[Service]
Environment=FLYDIGI_MODE=dinput
```

Then reload and restart the daemon:

```bash
sudo systemctl daemon-reload
sudo systemctl restart flydigid.service
```

This selects the protocol `flydigid` uses to configure the controller. It does not force the controller hardware into DInput/HID mode; use the controller's own mode switch/button sequence for that. Whether extra buttons such as `M1`-`M4` appear as unique Linux input buttons depends on the controller's HID reports in that mode.

## Basic Usage

Check CLI and daemon versions:

```bash
flydigictl version
```

Show connected controller information:

```bash
flydigictl info
```

Example:

```text
            Device : 105 (Vader 4 Pro)
     Protocol mode : dinput
Battery percentage : 100%
   Connection type : wireless
               CPU : wch (ch571)
```

Dump the parsed controller configuration:

```bash
flydigictl dump --no-color
```

Keep the daemon connected to the controller across commands:

```bash
flydigictl --persist-conn info
```

Force-close the controller connection after a command:

```bash
flydigictl --close-conn info
```

Print only the requested value for commands that support terse output:

```bash
flydigictl --terse joystick left deadzone
```

## Joystick Configuration

Read joystick deadzone:

```bash
flydigictl joystick left deadzone
flydigictl joystick right deadzone
```

Set joystick deadzone, from `0` to `100`:

```bash
flydigictl joystick left deadzone 8
flydigictl joystick right deadzone 8
```

## Button Mapping

Show the controller's current firmware button mappings:

```bash
flydigictl buttons
```

Map one controller button to another gamepad button:

```bash
flydigictl buttons map C A
flydigictl buttons map Z B
flydigictl buttons map M1 LB
flydigictl buttons map M2 RB
```

Reset a button to its native/unmapped value:

```bash
flydigictl buttons reset C
```

Recognized button names are `UP`, `RIGHT`, `DOWN`, `LEFT`, `A`, `B`, `SELECT`, `X`, `Y`, `START`, `LB`, `RB`, `LT`, `RT`, `THUMBL`, `THUMBR`, `C`, `Z`, `M1`, `M2`, `M3`, `M4`, `M5`, `M6`, `MENU`, `HOME`, and `BACK`. `L3` and `R3` are accepted aliases for `THUMBL` and `THUMBR`.

Button mapping is saved through the controller profile save path. In XInput mode, extra buttons such as `C`, `Z`, and `M1`-`M4` can be remapped to existing Xbox-style buttons, but Linux's `xpad` driver does not expose them as additional unique buttons.

## LED Configuration

Show the current LED mode. This reads the main config first so the active onboard profile ID is discovered before reading the separate LED config:

```bash
flydigictl leds
```

Turn LEDs off:

```bash
flydigictl leds off
```

Set a solid color:

```bash
flydigictl leds steady '#FF0000'
flydigictl leds steady red
flydigictl leds steady '#00F2FF' --brightness 0.5
```

Set the streamlined effect:

```bash
flydigictl leds streamlined --speed 0.5
flydigictl leds streamlined --speed 1.0 --brightness 0.75
```

Set the breathing effect:

```bash
flydigictl leds breathing '#BB00FF'
flydigictl leds breathing cyan --speed 0.5 --brightness 0.75
```

Set the gradient effect:

```bash
flydigictl leds gradient '#BB00FF' '#00F2FF'
flydigictl leds gradient purple cyan --speed 0.5 --brightness 0.75
```

Set the button feedback effect:

```bash
flydigictl leds feedback
flydigictl leds feedback --speed 0.5 --brightness 0.75
```

LED changes are written through the controller profile save path. This should make changes survive disconnects and controller power cycles on devices that persist LED settings in the active onboard profile.

Color arguments can be CSS color names from Go's `golang.org/x/image/colornames` package or hex colors in `#RRGGBB` format.

Known protocol LED modes:

- `Off`: supported for read/write
- `Streamlined`: supported for read/write
- `Breathing`: supported for read/write
- `Gradient`: supported for read/write
- `Feedback`: supported for read/write
- `Steady`: supported for read/write

If the controller is using an unknown mode, `flydigictl leds` prints `Unsupported LED mode`.

## Vader 4 Pro Notes

Vader 4 Pro reports device ID `105` and configuration version `1.3`. Its main configuration layout matches the Vader 3 Pro layout used by this tool.

Battery reporting on Vader 4 Pro is interpreted as a 0-4 charge level when the reported value is between `0` and `4`. Values above `4` and up to `100` are interpreted as direct percentages. If the value looks wrong, inspect the raw daemon log line:

```bash
sudo journalctl -u flydigid -n 30 --no-pager
```

Look for `got device info`; it includes the raw `battery` byte when debug logging is enabled.

## Game Input

`flydigictl` configures the controller over raw USB. Games do not use that raw USB connection; they need a Linux input device created by a kernel driver such as `xpad`.

Check whether the controller is bound to an input driver:

```bash
lsusb -t
ls -l /dev/input/by-id/ /dev/input/js*
```

For XInput mode, a working game input setup should show the controller with `Driver=xpad` and a joystick node such as `/dev/input/js0` or a by-id symlink like:

```text
/dev/input/by-id/usb-Flydigi_Flydigi_VADER4_Flydigi_VADER4-joystick
```

If `lsusb -t` shows the controller as `Driver=[none]`, reload `xpad` and replug the dongle:

```bash
sudo modprobe xpad
```

If this has to be repeated after every reconnect, check for an `xpad` blacklist. A blacklist prevents the kernel from auto-loading `xpad` when the controller appears:

```bash
ls /etc/modprobe.d/*xpad*.conf /usr/lib/modprobe.d/*xpad*.conf
```

Remove or edit any local file that contains `blacklist xpad`, then replug the controller. `flydigid` uses libusb auto-detach for XInput controllers, so it should not need to unload the global `xpad` module.

Test input events:

```bash
sudo evtest /dev/input/by-id/usb-Flydigi_Flydigi_VADER4_Flydigi_VADER4-event-joystick
jstest --event /dev/input/by-id/usb-Flydigi_Flydigi_VADER4_Flydigi_VADER4-joystick
```

On Arch Linux, install the test tools with:

```bash
sudo pacman -S evtest joystick
```

## Driver Conflicts

On Linux, kernel drivers can claim the controller before `flydigid` can communicate with it. For XInput controllers, `flydigid` asks libusb to detach and reattach the kernel driver while it is configuring the controller, so `xpad` usually does not need to be unloaded manually. Symptoms of a remaining driver conflict include:

- `flydigictl dump` returns `device doesn't respond`
- `flydigictl leds` fails while reading or writing
- daemon logs show the gamepad read loop exiting early

For XInput mode, the controller appears as:

```text
ID 045e:028e Microsoft Corp. Xbox360 Controller
```

Check loaded drivers:

```bash
lsmod | grep -E 'xpad|hid_xpadneo|uhid'
```

Temporarily unload `hid_xpadneo`:

```bash
sudo rmmod hid_xpadneo
```

If `hid_xpadneo` cannot be removed because the device is still bound, unbind only the matching HID device under `/sys/bus/hid/drivers/hid_xpadneo/`, then run `sudo rmmod hid_xpadneo` again.

These driver changes are session-only and reversible. Reload the module with:

```bash
sudo modprobe hid_xpadneo
```

They also revert on reboot unless you add persistent modprobe blacklist rules yourself.

## Troubleshooting

Controller is not detected:

```bash
flydigictl info
```

Check:

- The controller is in XInput or DInput mode.
- `lsusb` shows the controller.
- Conflicting kernel drivers are not holding the USB interface.
- `flydigid` can start and connect to DBus.

Inspect daemon logs:

```bash
sudo journalctl -u flydigid -f
```

Restart the daemon after replacing binaries:

```bash
sudo systemctl restart flydigid.service
```

## Uninstall

Arch Linux manual uninstall:

```bash
sudo systemctl stop flydigid.service || true
sudo rm -f /usr/bin/flydigictl /usr/bin/flydigid
sudo rm -f /usr/lib/systemd/system/flydigid.service
sudo rm -f /usr/share/dbus-1/system.d/flydigid.conf
sudo rm -f /usr/share/dbus-1/system-services/com.pipe01.flydigi.Gamepad.service
sudo systemctl daemon-reload
```

## Disclaimer

This program is provided as-is. Use it at your own risk.
