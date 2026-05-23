# Writing Custom Plugins for Wails Mobile

Wails Mobile uses a **Tri-Bridge Architecture** to connect the three layers of your application:
1. **Frontend (JavaScript/HTML/CSS)**: The user interface running in a WebView.
2. **Native Layer (Java/Android)**: Access to hardware and system APIs.
3. **Core Engine (Go)**: Business logic and application state.

This guide explains how to write a plugin that spans all three layers.

---

## 1. The Java Layer (Android)

Every plugin starts with a Java implementation that implements the `WailsPlugin` interface. This interface allows your plugin to hook into the Activity lifecycle and handle actions.

### Create the Plugin Class
Create a new package under `com.wailspackage.<yourplugin>` and implement `WailsPlugin`.

```java
package com.wailspackage.device;

import com.wailsplugin.WailsPlugin;
import android.os.Build;
import androidx.appcompat.app.AppCompatActivity;

public class DevicePlugin implements WailsPlugin {
    private AppCompatActivity mActivity;

    @Override
    public String getDomain() { return "device"; } // Unique namespace

    @Override
    public void onAttach(AppCompatActivity activity) {
        this.mActivity = activity;
    }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("getInfo".equals(action)) {
            return String.format("{\"model\":\"%s\", \"version\":\"%s\"}", 
                Build.MODEL, Build.VERSION.RELEASE);
        }
        return "{\"error\":\"Unknown action\"}";
    }

    @Override public void onActivityResult(int req, int res, Intent data) {}
    @Override public void onRequestPermissionsResult(int req, String[] perms, int[] res) {}
}
```

### Register the Plugin
In your `WailsWebViewActivity.java`, register the plugin in `onCreate`:

```java
registerPlugin(new DevicePlugin());
```

---

## 2. The Go Layer (Backend)

The Go layer acts as the orchestrator. It can call the Java layer and receive async callbacks.

### Implement the Go Plugin
Create a new directory under `wails/<yourplugin>`.

```go
package device

import (
	"github.com/its-ernest/wails-mobile/wails"
)

type DevicePlugin struct {
	app *wails.Application
}

func (p *DevicePlugin) Init(app *wails.Application) error {
	p.app = app
	return nil
}

// GetInfo calls the Java layer synchronously
func (p *DevicePlugin) GetInfo() string {
	return wails.HandleNativeAction("device:getInfo", "{}")
}
```

---

## 3. The Frontend Layer (JavaScript)

Finally, use the `WailsBind` object to interact with your plugin.

### Standard Synchronous Call
```javascript
async function getDeviceInfo() {
    // Call Go, which in turn calls Java
    const info = await WailsBind.callGo("DevicePlugin.GetInfo", "[]");
    console.log("Device Info:", JSON.parse(info));
}
```

### Direct Native Call (Bypassing Go)
```javascript
const info = WailsBind.callNative("device:getInfo", "{}");
```

---

## Best Practices

1. **Namespacing**: Always prefix your actions with your domain (e.g., `domain:action`).
2. **Asynchronous Events**: For long-running tasks (like taking a photo), return a "pending" status immediately from `handleAction` and use `app.Events.Emit` from Go to send the result back to JS when it's ready.
3. **JSON Consistency**: Always use JSON for payloads to ensure compatibility across the bridge.
