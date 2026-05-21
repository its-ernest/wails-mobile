# Wails Mobile

A small Go-based mobile app starter kit with a reusable runtime core and a hello-world Android template.

## Overview

- `wailsmobile/` - core Go runtime package for mobile bridged apps.
- `examples/hello-world/` - small sample app with embedded frontend, Go backend, and Android sync flow.
- `internal/templates/android/` - Android template app where generated AARs get dumped.
- `installer.sh` / `wailsm` - project starter that can be installed into your PATH.

## Install the CLI

Place the compiled `wailsm` binary into a folder on your `PATH`.

```bash
curl -L -o /usr/local/bin/wailsm https://github.com/its-ernest/wails-mobile/releases/latest/download/wailsm
chmod +x /usr/local/bin/wailsm
```

Now you can create a new app from anywhere:

```bash
wailsm my-app
```

That command clones the starter repo and sets up the minimal native project layout.

## Local development

If you want to work from source, clone this repo and use the example sync flow:

```bash
git clone https://github.com/its-ernest/wails-mobile.git
cd wails-mobile/examples/hello-world
chmod +x sync.sh
./sync.sh
```

The example sync script reads `config.ini`, builds the AAR with `gomobile bind`, and copies the generated library into `internal/templates/android/app/libs`.

## Configuring sync

Edit `examples/hello-world/config.ini` to change:

- `gomobile_target` - Android ABI(s) to build
- `androidapi` - Android API level for gomobile
- `aar_name` - output AAR file name
- `template_min_sdk` / `template_target_sdk` - values patched into the Android template

## Requirements

- Go 1.24+
- `gomobile` installed
- Android SDK + NDK installed for `gomobile bind`
- `git`, `curl`, `unzip` for the installer flow

## Notes

- The repo is designed so the hello-world example can compile against the core runtime in the same repository.
- Use `examples/hello-world/sync.sh` to keep the Android template libraries up to date.

