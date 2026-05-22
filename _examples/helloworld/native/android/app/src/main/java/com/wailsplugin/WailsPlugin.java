package com.wailsplugin;

import android.content.Intent;
import androidx.appcompat.app.AppCompatActivity;

public interface WailsPlugin {
    // The domain namespace this plugin handles (e.g., "permissions", "camera")
    String getDomain();
    
    // Core lifecycle attachment so the plugin can access context or trigger native prompts
    void onAttach(AppCompatActivity activity);
    
    // Handle actions routed from the Go/JS layer
    String handleAction(String action, String jsonArgsPayload);
    
    // Pass back activity results (crucial for permissions and intent pickers)
    void onActivityResult(int requestCode, int resultCode, Intent data);
    
    // Pass back permission result tokens
    void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults);
}