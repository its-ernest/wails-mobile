# Hello World Mobile Sample

A minimal mobile Go example packaged for Android using `gomobile bind`.

## Structure

- `enginego` — app startup and runtime registration
- `main.go` — backend service with `HelloService.SayHello`
- `frontend/` — embedded web UI assets
- `android/` — local Android build helper scripts

## Sync workflow

Create or update `config.ini` to control architecture, Android API, output path, and template SDK values.

```bash
cd examples/hello-world
chmod +x sync.sh
./sync.sh
```

The script will:

- build the example AAR with `gomobile bind`
- dump the generated `.aar` and source `.jar` into `internal/templates/android/app/libs`
- optionally patch `minSdk` and `targetSdk` in the Android template

## Build

```bash
cd examples/hello-world/android
chmod +x gen-aars.sh build-android.sh
./gen-aars.sh
```

This generates:

- `app/libs/wailsmobile.aar`
- `app/libs/helloworld.aar`

## Notes

- The generated AARs can be copied into `internal/templates/android/app/libs/` for Android Studio sync.
- The WebView frontend uses `WailsBind.callGo(...)` to invoke the Go backend.
