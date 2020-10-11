---
description: bluetooth headsets and linux don't like each other
title: bluetooth headset
---

### _bluetooth_ headset

So audio (in/out)
Sony WH-1000MX3 and Intel® Dual Band Wireless-AC 8265 on aXPS13
running Arch Linux, how hard can it be?

#### _basics_

_bluez_ manages the bluetooth connection part, but will fail if it doesn't have access to the right profile

#### _install_

- **bluez**: bluetooth
- **bluez-utils**: actually control bluetooth
- **pulseaudio**: do you really want to make your life harder by trying ALSA?
- **pulueaudio-bluetooth**: make it talk with bluez
- **pulsemixer**: (optional) ncurses interface to control pulseaudio,
  alternatively pamixer for pure cli, pavucontrol for gui

#### _post_ install

- `systemctl --user enable --now pulseaudio` pulseaudio has socket activation,
  but it doesn't work if you try to connect to a bluetooth headset before something has started pulseaudio,
  so enable it to always make it available.
- `bluetoothctl power on` turn on bluetooth
  (? or `AutoEnable=true` in `/etc/bluetooth/main.conf` but not sure how it interacts with pulseaudio being user service)
- `bluetoothctl agent on` only needs to be done once
- `bluetoothctl scan on` find bluetooth things in range, only needs to be done once
- `bluetoothctl pair xx:xx:xx:xx:xx:xx` pair with device, only needs to be done once
- `bluetoothctl connect xx:xx:xx:xx:xx:xx` connect to device
- `bluetoothctl trust xx:xx:xx:xx:xx:xx` autoconnect to device on power on

#### _sound_ control

- **Output**: where sound goes, can set a default device
- **Input**: mic control, `Monitor of ...` is to listen to what the mic picked up
- **Cards**: sound profile control, _A2DP_ for high quality output or _HSP/HFP_ for tin can quality mic+output

#### _todo_

- `AutoEnable=true` in `/etc/bluetooth/main.conf`
- `load-module module-switch-on-connect` in `.config/pulse/default.pa`

+more pulseaudio modules
