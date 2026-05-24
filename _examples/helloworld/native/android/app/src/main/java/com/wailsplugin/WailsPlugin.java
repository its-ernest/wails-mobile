package com.wailsplugin;

import android.content.Context;
import android.content.Intent;

/**
 * WailsPlugin is the base interface for all native Android extensions.
 */
public interface WailsPlugin {
    
    String getDomain();
    
    /**
     * Called when the plugin is registered or when an activity becomes active.
     * @param context Can be an Application context or an AppCompatActivity context.
     */
    void onAttach(Context context);
    
    String handleAction(String action, String jsonArgsPayload);
    
    void onActivityResult(int requestCode, int resultCode, Intent data);
    
    void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults);
}
