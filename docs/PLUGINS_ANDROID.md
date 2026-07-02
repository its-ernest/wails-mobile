# Android Plugin Development for Wails Mobile

This guide covers the technical details of implementing the native side of a Wails Mobile plugin on Android.

## The `WailsPlugin` Interface

All Android plugins must implement the `com.wailsplugin.WailsPlugin` interface:

```java
public interface WailsPlugin {
    // Unique domain for routing (e.g., "camera", "permission")
    String getDomain();
    
    // Lifecycle hook for context attachment
    void onAttach(Context context);
    
    // Primary execution entry point from Go
    String handleAction(String action, String jsonArgsPayload);
    
    // Lifecycle hooks for Android results
    void onActivityResult(int req, int res, Intent data);
    void onRequestPermissionsResult(int req, String[] perms, int[] res);
}
```

---

## Background-Safe Design

Android apps can be woken up in the background (e.g., by `WorkManager`). In these cases, there is **no Activity**. Your plugin must be designed to handle this.

### Context Awareness
`onAttach(Context context)` is called twice:
1.  **At Startup**: Called with the `Application` context. Use this for system services (e.g `NotificationManager`, `Log`).
2.  **On UI Open**: Called with the `AppCompatActivity` context. Use this for UI tasks (e.g `ActivityCompat.requestPermissions`, `startActivityForResult`).

### Example: Context Handling
```java
public class MyPlugin implements WailsPlugin {
    private Context mContext;
    private AppCompatActivity mActivity;

    @Override
    public void onAttach(Context context) {
        this.mContext = context; // Always available
        if (context instanceof AppCompatActivity) {
            this.mActivity = (AppCompatActivity) context; // Available only in foreground
        }
    }

    @Override
    public String handleAction(String action, String args) {
        if ("doUiTask".equals(action)) {
            if (mActivity == null) return "{\"error\":\"App in background\"}";
            // Perform Activity-specific task
        }
        return "{\"status\":\"ok\"}";
    }
}
```

---

## Registration

### 1. The Global scope
Register your plugin in `WailsApplication.java` to ensure it is available for background workers
Your plugin's end-developers will register the plugin for their app manually:

```java
// Inside registerPlugins() in WailsApplication.java
private void registerPlugins() {
    registerPlugin(new MyPlugin());
}
```

### 2. The UI Scope
Ensure your plugin is re-attached to the Activity in `WailsWebViewActivity.java`:
This is usually handled or done automatically in Wails Mobile. Only add it if it doesn't exist. Standard Wails Mobile project comes with it.
```java
// Inside onCreate
WailsApplication app = (WailsApplication) getApplication();
for (WailsPlugin plugin : app.getPlugins().values()) {
    plugin.onAttach(this);
}
```

---

## Talking back to Go

To send data from Java back to Go asynchronously:

```java
// Wrap your arguments in a JSON array []
String payload = "[{\"result\": \"success\"}]";

// handleNativeAction routes directly to Go's RegisterNativeMethod handlers
Wailsmobile.handleNativeAction("myplugin:callback", payload);
```

---

## Tips for Android
*   **Permissions**: Use the existing `PermissionsPlugin` to check/request neccessary permissions before your plugin executes.
*   **Main Thread**: `handleAction` is often called on a background bridge thread. If you need to touch the UI, use `new Handler(Looper.getMainLooper()).post(...)`.
*   **JSON Parsing**: Use `org.json.JSONObject` to parse `jsonArgsPayload`.

---

## 🛠️ Background Services & Daemons

If your plugin requires a persistent background presence, you should implement a `Service`.

1.  **Foreground Service**: Required for persistent background logic on modern Android (Oreo+). You must show a notification.
2.  **Manifest Registration**: All services must be declared in `AndroidManifest.xml`.
3.  **Permissions**: Request `FOREGROUND_SERVICE` and specific types (e.g., `specialUse` for custom Go logic) in the manifest.
