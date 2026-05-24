# Wails Mobile (Stable Android Support)

**Automation scripts currently works for Linux systems**

Wails Mobile is a port of Go v3 implementation to support mobile devies
Gradle support: `Gradle 9.2.1` by default. You can change the version to your Gradle version

## Overview

- `wails/` - core Go runtime package for mobile bridged apps.
- `_examples/helloworld/` - small sample app with embedded frontend, Go backend, and Android sync flow.
- `native/android/` - Android app template app where generated AARs get dumped.
- `wailsm.sh` / `wailsm` - project starter and CLI automated helper that can be installed into your PATH.

---

## Install the CLI

To install  and build apps with `wails-mobile`, use the CLI helper by running the command below:
Download size is below: `700 kb`
Use the command below:

```bash
sudo curl -L -H "Cache-Control: no-cache" -o /usr/local/bin/wailsm https://github.com/its-ernest/wails-mobile/releases/latest/download/wailsm && \
sudo chmod +x /usr/local/bin/wailsm
```

Now you can create a new app from anywhere:

```bash
wailsm --new my-new-app
```

That command downloads a starter pre-configured Wails Mobile project files into `my-app/` directory, ready to be edited and compiled.

The pre-configured project can be extended(adding more code), 
when it's time to build, open `native/android/` folder in Android Studio and click `Build`.

--- 

## Requirements

- Go 1.24+
- Android SDK + NDK installed for `gomobile bind`
- `git`, `curl`, `unzip` for the installer flow

---

## WailsMobile App Structure

- `engine.go` — app startup and runtime registration, not meant to be tampered unless you are sure
- `main.go` — backend service with `AppService.SayHello`
- `frontend/` — embedded web UI assets
- `native/` — this folder contains the standard Android Studio or XCode app project

## Refresh or Sync workflow

After editing or write your Go code or frontend, you must refresh or sync first to produce native bindings.

You can optionally update `android.ini` to control architecture, Android API, output path, and template SDK values. By default, this works for Android 5 up to latest Android 17.

```bash
wailsm --refresh android # or ios
```

The script will:

- build the example AAR with `gomobile bind`
- dump the generated `.aar` and source `.jar` into `native/android/app/libs`


## Building the App

Once again: don't forget after editing or write your Go code, you must refresh or sync first to produce native bindings.

Now open the project in Android Studio or XCode.

Where the resulting projects are found:
 - Final Android Studio project: `native/android/`
 - Final XCode: `native/...`


## Adding external Plugins

External plugins are just external Go packages, with the directory structured in such a way Wails Mobile CLI can realize it contains some native platform code. 

Sample Plugin structure:

```bash
somePlugin/
    - android/
        - com/.... #example package name
    - ios/
        - ...#support coming soon
    
    - go.mod
    - other go package files
```

You can add a plugin with:
```bash
wailsm --add github.com/<handle>/<plugin> # uses `go get...` under the hood
```

## Removing packages
You can remove a plugin from Wails mobile with:
```bash
wailsm --remove github.com/<handle>/<plugin>
```

## Writing custom plugins or Native code

You can design an external plugin to be integrated into Wails Mobile application, or even write some app-specific native Java or Swift code and call it from Go. 

Instructions on achieving this goal is wired in `[PLUGINS.md](PLUGINS.md)`

---

## To uninstall

```bash
sudo rm -rf /usr/local/bin/wailsm
```


## Notes

- The WebView frontend uses `WailsBind.callGo(...)` to invoke the Go backend.
- Currently, `wails-mobile` is stable enough for Android builds. If you care to get support for iOS sooner, contribute by translating the Java project into Swift project. 
- The core bridge in Go is cross-platform(`Android` and `iOS`). No writing of JNI. No writing of *C* code and headers. 
