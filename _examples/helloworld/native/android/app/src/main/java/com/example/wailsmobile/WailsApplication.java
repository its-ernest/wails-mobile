package com.example.wailsmobile;

import android.app.Application;
import com.wailsplugin.WailsPlugin;
import com.wailspackage.biometric.BiometricPlugin;
import com.wailspackage.devicestate.DeviceStatePlugin;
import com.wailspackage.filepicker.FilePickerPlugin;
import com.wailspackage.logger.LoggerPlugin;
import com.wailspackage.notifications.NotificationPlugin;
import com.wailspackage.permissions.PermissionsPlugin;
import com.wailspackage.workmanager.WorkManagerPlugin;
import java.util.HashMap;
import java.util.Map;
import wailsmobile.Wailsmobile;

public class WailsApplication extends Application {
    private final Map<String, WailsPlugin> mPlugins = new HashMap<>();

    @Override
    public void onCreate() {
        super.onCreate();
        
        // Start Go engine
        Wailsmobile.startApplication();

        // Initialize and Register global plugins
        registerPlugin(new PermissionsPlugin());
        registerPlugin(new WorkManagerPlugin());
        registerPlugin(new NotificationPlugin());
        registerPlugin(new LoggerPlugin());
        registerPlugin(new DeviceStatePlugin());
        registerPlugin(new BiometricPlugin());
        registerPlugin(new FilePickerPlugin());

        // Register the global handler for Go-to-Native calls
        Wailsmobile.setNativeCallHandler(new wailsmobile.NativeCallHandler() {
            @Override
            public String onNativeCall(String method, String args) {
                if (method.contains(":")) {
                    String[] parts = method.split(":", 2);
                    String domain = parts[0];
                    String action = parts[1];

                    WailsPlugin plugin = mPlugins.get(domain);
                    if (plugin != null) {
                        return plugin.handleAction(action, args);
                    }
                }
                return "{\"error\":\"Plugin domain not found\"}";
            }
        });
    }

    private void registerPlugin(WailsPlugin plugin) {
        // Prime the plugin with the Application context for background tasks
        plugin.onAttach(this);
        mPlugins.put(plugin.getDomain(), plugin);
    }

    public Map<String, WailsPlugin> getPlugins() {
        return mPlugins;
    }
}
