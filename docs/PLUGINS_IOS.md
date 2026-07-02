# iOS Plugin Development for Wails Mobile (WIP)

This guide covers the technical details of implementing the native side of a Wails Mobile plugin on iOS.

> **Note**: iOS support is currently experimental and uses a Docker-based cross-compilation toolchain.

## 📋 The Swift/Objective-C Interface

Similar to Android, all iOS plugins must implement a common interface (Protocol in Swift).

```swift
protocol WailsPlugin {
    func getDomain() -> String
    func onAttach(container: WailsContainer)
    func handleAction(action: String, jsonArgs: String) -> String
}
```

(Detailed interface definition coming as we solidify the `xtools` integration.)

---

## 🏗️ Docker-Based Toolchain

Wails Mobile uses `docker-compose` to manage the iOS build environment on non-macOS systems.

1.  **Go Bindings**: Generated inside the container using `gomobile bind`.
2.  **Swift Compilation**: Managed via `xtools` and the Swift Package Manager inside the container.
3.  **Deployment**: Bridged via `usbmuxd` from the host to the container.

---

## 📡 Talking back to Go

Async communication from Swift back to Go follows the same pattern as Android:

```swift
let payload = "[{\"result\": \"success\"}]"
WailsmobileHandleNativeAction("myplugin:callback", payload)
```
