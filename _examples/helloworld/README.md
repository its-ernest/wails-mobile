# Hello World Mobile Sample

A minimal mobile Go example packaged for Android using `gomobile bind`.

## Structure

- `engine.go` — app startup and runtime registration, not meant to be tampered unless you are sure
- `main.go` — backend service with `AppService.SayHello`
- `frontend/` — embedded web UI assets
- `android/` — local Android build helper scripts

## Refresh or Sync workflow

After editing or write your Go code, you must refresh or sync first to produce native bindings.

Update `android.ini` to control architecture, Android API, output path, and template SDK values.

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

## Notes

- The WebView frontend uses `WailsBind.callGo(...)` to invoke the Go backend.
- Currently, `wails-mobile` is stable enough for Android builds. If you care to get support for iOS sooner, contribute by translating the Java project into Swift project. 
- The core bridge in Go is cross-platform(`Android` and `iOS`). No writing of JNI. No writing of *C* code and headers. 
