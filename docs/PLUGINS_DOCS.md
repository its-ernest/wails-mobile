# Wails Mobile Plugin Development Guide

Wails Mobile utilizes a **Go-Centric Tri-Bridge Architecture**. This design ensures that business logic remains in Go while providing seamless access to native platform features and web-based user interfaces.

## The Architecture

1.  **Frontend (UI)**: Dispatches user actions to Go.
2.  **Go Engine (Logic & Orchestration)**: The "Brain" of the application. It processes requests, manages state, and instructs the Native Layer when platform-specific hardware or APIs are required.
3.  **Native Layer (Platform Features)**: Fulfills requests from Go (e.g., Camera, Notifications, WorkManager). It is **Context-Aware**, supporting both foreground (Activity) and background (Service/Worker) operations.

---

## Developing a Plugin

A complete plugin typically consists of three parts:

### 1. The Go Wrapper (`wails/myplugin`)
The Go layer provides the primary API for both the Frontend and other Go services.

```go
package myplugin

import "github.com/its-ernest/wails-mobile/wails"

type MyPlugin struct {
    app *wails.Application
}

//Init attaches Application context in Go to your plugin
func (p *MyPlugin) Init(app *wails.Application) error {
    p.app = app
    return nil
}

// PerformAction can be called by JS or Go. It orchestrates the native call
func (p *MyPlugin) PerformAction(data string) string {
    // Go calls the Native Platform bridge
    // the key "myplugin:action" would be implemented in native Java or Swift
    // the end-dev has to never write JNI
    return wails.CallNativePlatform("myplugin:action", data)
}
```

### 2. The Native Implementation
The native implementation handles the actual platform API calls. It must implement the `WailsPlugin` interface.

*   [**Android Implementation Guide (Java/Kotlin)**](./PLUGINS_ANDROID.md)
*   [**iOS Implementation Guide (Swift/Obj-C)**](./PLUGINS_IOS.md) (WIP - Experimental)

### 3. The Frontend Integration (`app.js`)
The frontend interacts with the plugin **strictly through the Go layer**.

```js
// Call the Go wrapper
// similar to how you'd call Go methods
const result = await Wails.CallGo('MyPlugin.PerformAction', "some-data");
```

---

## Bidirectional Communication

### Synchronous (Request-Response)
Standard calls like `wails.CallNativePlatform` return a string response immediately.

### Asynchronous (Events)
For long-running tasks, the Native Layer should trigger a "Native Action" back to Go, which Go then emits as a Wails Event:

1.  **Native**: Calls `Wailsmobile.handleNativeAction("myplugin:finished", json)`.
2.  **Go**: In plugin's `Init`, registers a handler:
    ```go
    app.RegisterNativeMethod("myplugin:finished", func(args []json.RawMessage) (interface{}, error) {
        app.Events.Emit("plugin-event", args[0])
        return nil, nil
    })
    ```
3.  **Frontend**: Listens for the event:
    ```js
    Wails.on("plugin-event", (data) => { ... });
    ```

---

## Best Practices

*   **Go-First Logic**: Never put business logic in the Native layer. Use Native only for API access.
*   **JSON Everywhere**: Always use JSON strings for payloads between Go and Native to ensure cross-platform compatibility.
*   **Context Safety**: Ensure your native code is background-safe (see platform-specific guides).
*   **Namespacing**: Prefix all native actions with your plugin domain (e.g., `logger:log`, `camera:capture`).