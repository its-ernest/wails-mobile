package com.wailsplugin;

import android.content.Intent;
import androidx.appcompat.app.AppCompatActivity;

/**
 * WailsPlugin is the base interface for all native Android extensions.
 * Plugins allow the Go and JavaScript layers to access native platform features
 * like cameras, sensors, and system permissions.
 */
public interface WailsPlugin {
    
    /**
     * Returns the unique domain/namespace for this plugin.
     * All calls from JavaScript starting with "domain:" will be routed here.
     * @return The plugin domain (e.g., "permissions")
     */
    String getDomain();
    
    /**
     * Called when the plugin is registered with the main Wails activity.
     * Provides access to the activity context for UI operations and system services.
     * @param activity The host WailsWebViewActivity
     */
    void onAttach(AppCompatActivity activity);
    
    /**
     * Handles an action request routed from the JavaScript or Go layer.
     * @param action The specific action to perform (e.g., "request")
     * @param jsonArgsPayload The JSON-formatted arguments for the action
     * @return A JSON-formatted string response
     */
    String handleAction(String action, String jsonArgsPayload);
    
    /**
     * Forwards the activity's onActivityResult callback to the plugin.
     * Essential for plugins that launch external intents (e.g., File Pickers).
     */
    void onActivityResult(int requestCode, int resultCode, Intent data);
    
    /**
     * Forwards the activity's onRequestPermissionsResult callback to the plugin.
     * Essential for plugins that request Android runtime permissions.
     */
    void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults);
}
