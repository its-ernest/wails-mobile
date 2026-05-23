package com.example.wailsmobile;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
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
import wailsmobile.Wailsmobile;

/**
 * WailsWebViewActivity is the primary container for the Wails Mobile application.
 * It manages the WebView lifecycle, handles asset loading via the Go backend,
 * and maintains the plugin registry for native extensions.
 */
public class WailsWebViewActivity extends AppCompatActivity {

    private WebView mWebView;
    private final Map<String, WailsPlugin> mPlugins = new HashMap<>();
    private final Handler mHandler = new Handler(Looper.getMainLooper());
    private boolean mIsPolling = true;

    /**
     * Registers a native plugin with the activity.
     * The plugin will be attached to this activity and its domain will be used for routing.
     * @param plugin The plugin instance to register
     */
    public void registerPlugin(WailsPlugin plugin) {
        plugin.onAttach(this);
        mPlugins.put(plugin.getDomain(), plugin);
    }

    @SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        // Initialize the Go backend
        String initResult = Wailsmobile.startApplication();
        Log.d("Wailsmobile", "Go backend init result: " + initResult);

        // Register Java handler for calls originating from Go
        // The interface is generated in the wailsmobile package
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

        mWebView = new WebView(this);

        // Force hardware acceleration for WebView rendering on Android versions that support it
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.KITKAT) {
            getWindow().setFlags(
                android.view.WindowManager.LayoutParams.FLAG_HARDWARE_ACCELERATED,
                android.view.WindowManager.LayoutParams.FLAG_HARDWARE_ACCELERATED
            );
        }

        // Enable algorithmic darkening (Dark Mode support)
        if (androidx.webkit.WebViewFeature.isFeatureSupported(androidx.webkit.WebViewFeature.ALGORITHMIC_DARKENING)) {
            if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.TIRAMISU) {
                mWebView.getSettings().setAlgorithmicDarkeningAllowed(true);
            }
        }

        setContentView(mWebView);

        // Configure WebView settings for Wails
        mWebView.getSettings().setJavaScriptEnabled(true);
        mWebView.getSettings().setDomStorageEnabled(true);

        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.LOLLIPOP) {
            mWebView.getSettings().setMixedContentMode(android.webkit.WebSettings.MIXED_CONTENT_ALWAYS_ALLOW);
        }

        mWebView.getSettings().setAllowFileAccess(true);
        mWebView.getSettings().setAllowContentAccess(true);
        mWebView.getSettings().setAllowFileAccessFromFileURLs(true);
        mWebView.getSettings().setAllowUniversalAccessFromFileURLs(true);
        mWebView.getSettings().setUseWideViewPort(true);
        mWebView.getSettings().setLoadWithOverviewMode(true);

        // Set up the JavaScript-to-Native bridge
        mWebView.addJavascriptInterface(new Object() {
            /**
             * Directly calls a Go method bound to the application.
             * The Go method is responsible for any further native calls.
             */
            @JavascriptInterface
            public String callGo(String methodKey, String jsonArgsPayload) {
                return Wailsmobile.handleMessageFromFrontend(methodKey, jsonArgsPayload);
            }
        }, "WailsBind");

        // Register default plugins
        registerPlugin(new PermissionsPlugin());

        mWebView.setWebViewClient(new WebViewClient() {
            /**
             * Intercepts URL requests to serve assets directly from the Go embedded filesystem.
             */
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                String urlStr = request.getUrl().toString();
                
                // Strip out the artificial base URL domain to isolate the relative asset key
                String assetKey = urlStr.replace("https://wails.local/", "");
                
                // Remove trailing query hashes or params
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
                    mimeType = "text/plain";
                }

                return new WebResourceResponse(
                        mimeType,
                        "UTF-8",
                        new ByteArrayInputStream(fileBytes)
                );
            }

            @Override
            public void onPageStarted(WebView view, String url, android.graphics.Bitmap favicon) {
                super.onPageStarted(view, url, favicon);
                // Inject the event system as early as possible. 
                // We use window.WailsEvents as a backup if WailsBind is not extensible.
                String js = "if (!window.WailsEvents) {" +
                        "  window.WailsEvents = {" +
                        "    listeners: {}," +
                        "    on: function(name, cb) { " +
                        "      if(!this.listeners[name]) this.listeners[name] = []; " +
                        "      this.listeners[name].push(cb); " +
                        "    }," +
                        "    dispatch: function(obj) { " +
                        "      var name = obj.name; var data = obj.data; " +
                        "      if(this.listeners[name]) { " +
                        "        this.listeners[name].forEach(function(cb) { try { cb(data); } catch(e) { console.error(e); } }); " +
                        "      } " +
                        "    }" +
                        "  };" +
                        "  if (window.WailsBind) {" +
                        "    window.WailsBind.on = window.WailsEvents.on.bind(window.WailsEvents);" +
                        "    window.WailsBind.dispatch = window.WailsEvents.dispatch.bind(window.WailsEvents);" +
                        "  }" +
                        "}";
                view.evaluateJavascript(js, null);
            }
        });

        // Start the application by loading the local entry point
        mWebView.loadUrl("https://wails.local/");

        startEventPolling();
    }

    /**
     * Starts a background thread to poll for events from the Go backend.
     * Events are dispatched to the WebView via JavaScript.
     */
    private void startEventPolling() {
        new Thread(() -> {
            while (mIsPolling) {
                String eventJson = Wailsmobile.pollNativeEvent();
                if (eventJson != null && !eventJson.isEmpty()) {
                    mHandler.post(() -> {
                        // Inject JS to dispatch the event. Try WailsBind.dispatch first, fallback to WailsEvents.dispatch.
                        String script = "if(window.WailsBind && window.WailsBind.dispatch) { " +
                                "window.WailsBind.dispatch(" + eventJson + "); " +
                                "} else if(window.WailsEvents && window.WailsEvents.dispatch) { " +
                                "window.WailsEvents.dispatch(" + eventJson + "); " +
                                "}";
                        mWebView.evaluateJavascript(script, null);
                    });
                }
                try {
                    Thread.sleep(100); // Poll every 100ms
                } catch (InterruptedException e) {
                    break;
                }
            }
        }).start();
    }

    @Override
    protected void onPause() {
        super.onPause();
        if (mWebView != null) {
            mWebView.onPause();
        }
        Wailsmobile.handleNativeAction("lifecycle:pause", "");
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (mWebView != null) {
            mWebView.onResume();
        }
        Wailsmobile.handleNativeAction("lifecycle:resume", "");
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        // Forward permission results to all registered plugins
        for (WailsPlugin plugin : mPlugins.values()) {
            plugin.onRequestPermissionsResult(requestCode, permissions, grantResults);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        // Forward activity results to all registered plugins
        for (WailsPlugin plugin : mPlugins.values()) {
            plugin.onActivityResult(requestCode, resultCode, data);
        }
    }

    @Override
    protected void onDestroy() {
        mIsPolling = false;
        // Gracefully stop the Go runtime engine before closing the activity
        Wailsmobile.handleNativeAction("lifecycle:destroy", "");
        if (mWebView != null) {
            mWebView.destroy();
        }
        super.onDestroy();
    }
}
