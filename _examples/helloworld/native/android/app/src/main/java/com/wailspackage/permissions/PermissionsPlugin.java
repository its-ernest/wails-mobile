package com.wailspackage.permissions;

import android.content.Intent;
import android.content.pm.PackageManager;
import android.util.Log;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import com.wailsplugin.WailsPlugin;

import org.json.JSONException;
import org.json.JSONObject;

import wailsmobile.Wailsmobile;

public class PermissionsPlugin implements WailsPlugin {
    private AppCompatActivity mActivity;
    private static final int PERMISSION_REQ_CODE = 9911;

    @Override
    public String getDomain() { return "permissions"; }

    @Override
    public void onAttach(AppCompatActivity activity) { this.mActivity = activity; }

    @Override
    public String handleAction(String action, String jsonArgsPayload) {
        if ("check".equals(action)) {
            // Parse jsonArgsPayload to find out which permission to check
            String perm = parsePermissionFromJson(jsonArgsPayload);
            int result = ContextCompat.checkSelfPermission(mActivity, perm);
            return result == PackageManager.PERMISSION_GRANTED ? "{\"status\":\"granted\"}" : "{\"status\":\"denied\"}";
        }

        if ("request".equals(action)) {
            String perm = parsePermissionFromJson(jsonArgsPayload);
            ActivityCompat.requestPermissions(mActivity, new String[]{perm}, PERMISSION_REQ_CODE);
            return "{\"status\":\"requested\"}";
        }

        return "{\"error\":\"Unknown action\"}";
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, String[] permissions, int[] grantResults) {
        if (requestCode == PERMISSION_REQ_CODE && grantResults.length > 0) {
            boolean granted = grantResults[0] == PackageManager.PERMISSION_GRANTED;
            String permission = permissions[0];

            // Native callback: async pass the permission status straight back into Go!
            // We use a domain-prefixed event name so Go knows how to route it.
            JSONObject result = new JSONObject();
            try {
                result.put("permission", permission);
                result.put("granted", granted);
                WailsMobile.handleMessageFromFrontend("permissions:result", result.toString());
            } catch (JSONException e) {
                Log.e("PermissionsPlugin", "Error creating result JSON", e);
            }
        }
    }

    @Override public void onActivityResult(int r, int rc, Intent d) {}

    private String parsePermissionFromJson(String json) {
        try {
            JSONObject obj = new JSONObject(json);
            return obj.optString("permission", android.Manifest.permission.CAMERA);
        } catch (JSONException e) {
            return android.Manifest.permission.CAMERA;
        }
    }
}