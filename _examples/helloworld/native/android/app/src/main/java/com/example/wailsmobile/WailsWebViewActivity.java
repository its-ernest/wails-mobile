package com.example.wailsmobile;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.webkit.JavascriptInterface;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebView;
import android.webkit.WebViewClient;

import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.Map;

import com.wailsplugin.WailsPlugin;
import com.wailspackage.permissions.PermissionsPlugin;
import wailsmobile.wailsmobile.Wailsmobile;

public class WailsWebViewActivity extends AppCompatActivity {

    private WebView mWebView;
    private final Map<String, WailsPlugin> mPlugins = new HashMap<>();

    public void registerPlugin(WailsPlugin plugin) {
        plugin.onAttach(this);
        mPlugins.put(plugin.getDomain(), plugin);
    }

    @SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        String initResult = Wailsmobile.startApplication();
        Log.d("Wailsmobile", "Go backend init result: " + initResult);

        mWebView = new WebView(this);

        // Force hardware acceleration for WebView rendering on Android versions that support it
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.KITKAT) {
            // force hardware acceleration at the window layer
            getWindow().setFlags(
                android.view.WindowManager.LayoutParams.FLAG_HARDWARE_ACCELERATED,
                android.view.WindowManager.LayoutParams.FLAG_HARDWARE_ACCELERATED
            );
        }

        if (androidx.webkit.WebViewFeature.isFeatureSupported(androidx.webkit.WebViewFeature.ALGORITHMIC_DARKENING)) {
            if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.TIRAMISU) {
                mWebView.getSettings().setAlgorithmicDarkeningAllowed(true);
            }
        }

        setContentView(mWebView);

        mWebView.getSettings().setJavaScriptEnabled(true);
        mWebView.getSettings().setDomStorageEnabled(true);

        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.LOLLIPOP) {
            mWebView.getSettings().setMixedContentMode(android.webkit.WebSettings.MIXED_CONTENT_ALWAYS_ALLOW);
        }

        mWebView.getSettings().setAllowFileAccess(true);
        mWebView.getSettings().setAllowContentAccess(true);
        // Allow resources loaded from local contexts to access other local endpoints
        mWebView.getSettings().setAllowFileAccessFromFileURLs(true);
        mWebView.getSettings().setAllowUniversalAccessFromFileURLs(true);

        mWebView.getSettings().setUseWideViewPort(true);
        mWebView.getSettings().setLoadWithOverviewMode(true);

        mWebView.addJavascriptInterface(new Object() {
            @JavascriptInterface
            public String callGo(String methodKey, String jsonArgsPayload) {
                return Wailsmobile.handleMessageFromFrontend(methodKey, jsonArgsPayload);
            }

            @JavascriptInterface
            public String callNative(String methodKey, String jsonArgsPayload) {
                // Expecting methodKey format -> "domain:action" (e.g., "permissions:request")
                if (methodKey.contains(":")) {
                    String[] parts = methodKey.split(":", 2);
                    String domain = parts[0];
                    String action = parts[1];

                    WailsPlugin plugin = mPlugins.get(domain);
                    if (plugin != null) {
                        return plugin.handleAction(action, jsonArgsPayload);
                    }
                }
                // Fallback to the default bridge if no plugin matches
                return Wailsmobile.handleNativeAction(methodKey, jsonArgsPayload);
            }
        }, "WailsBind");

        // Register default plugins
        registerPlugin(new PermissionsPlugin());

        mWebView.setWebViewClient(new WebViewClient() {
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                String urlStr = request.getUrl().toString();
                
                // Strip out the artificial base URL domain to isolate the relative asset key
                String assetKey = urlStr.replace("https://wails.local/", "");
                
                // Remove trailing query hashes or params (e.g. "main.js?v=123" -> "main.js")
                if (assetKey.contains("?")) {
                    assetKey = assetKey.split("\\?")[0];
                }
                if (assetKey.contains("#")) {
                    assetKey = assetKey.split("#")[0];
                }

                if (assetKey.isEmpty() || assetKey.equals("/")) {
                    assetKey = "index.html";
                }

                byte[] fileBytes = Wailsmobile.requestAssetBytes(assetKey);
                String mimeType = Wailsmobile.requestAssetMime(assetKey);
                if (mimeType == null || mimeType.isEmpty()) {
                    mimeType = "text/plain"; // WebViews handle missing mimes poorly
                }

                return new WebResourceResponse(
                        mimeType,
                        "UTF-8",
                        new ByteArrayInputStream(fileBytes)
                );
            }
        });

        String html = new String(Wailsmobile.requestAssetBytes("index.html"), StandardCharsets.UTF_8);
        mWebView.loadUrl("https://wails.local/");
    }

    @Override
    protected void onPause() {
        super.onPause();
        if (mWebView != null) {
            mWebView.onPause(); // Pauses JavaScript timers and flash animations
        }
        Wailsmobile.handleNativeAction("lifecycle:pause", ""); // Notify Go
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (mWebView != null) {
            mWebView.onResume(); // Resumes JavaScript execution
        }
        Wailsmobile.handleNativeAction("lifecycle:resume", ""); // Notify Go
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        for (WailsPlugin plugin : mPlugins.values()) {
            plugin.onRequestPermissionsResult(requestCode, permissions, grantResults);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        for (WailsPlugin plugin : mPlugins.values()) {
            plugin.onActivityResult(requestCode, resultCode, data);
        }
    }

    @Override
    protected void onDestroy() {
        // Gracefully stop the Go runtime engine before closing the activity
        Wailsmobile.handleNativeAction("lifecycle:destroy", "");
        if (mWebView != null) {
            mWebView.destroy();
        }
        super.onDestroy();
    }
}
