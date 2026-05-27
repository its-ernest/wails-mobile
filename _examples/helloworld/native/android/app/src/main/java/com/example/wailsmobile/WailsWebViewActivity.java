package com.example.wailsmobile;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import android.webkit.JavascriptInterface;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;
import com.wailsplugin.WailsPlugin;
import java.io.ByteArrayInputStream;
import java.util.Map;
import wailsmobile.Wailsmobile;

public class WailsWebViewActivity extends AppCompatActivity {

    private WebView mWebView;
    private final Handler mHandler = new Handler(Looper.getMainLooper());
    private boolean mIsPolling = true;

    @SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        // Attach this activity context to all plugins so they can perform UI actions/permissions
        WailsApplication app = (WailsApplication) getApplication();
        for (WailsPlugin plugin : app.getPlugins().values()) {
            plugin.onAttach(this);
        }

        mWebView = new WebView(this);

        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.KITKAT) {
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
        mWebView.getSettings().setAllowFileAccessFromFileURLs(true);
        mWebView.getSettings().setAllowUniversalAccessFromFileURLs(true);
        mWebView.getSettings().setUseWideViewPort(true);
        mWebView.getSettings().setLoadWithOverviewMode(true);

        mWebView.addJavascriptInterface(new Object() {
            @JavascriptInterface
            public String callGo(String methodKey, String jsonArgsPayload) {
                return Wailsmobile.handleMessageFromFrontend(methodKey, jsonArgsPayload);
            }
        }, "WailsBind");

        mWebView.setWebViewClient(new WebViewClient() {
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                String urlStr = request.getUrl().toString();
                String assetKey = urlStr.replace("https://wails.local/", "");
                
                if (assetKey.contains("?")) assetKey = assetKey.split("\\?")[0];
                if (assetKey.contains("#")) assetKey = assetKey.split("#")[0];
                if (assetKey.isEmpty() || assetKey.equals("/")) assetKey = "index.html";

                byte[] fileBytes = Wailsmobile.requestAssetBytes(assetKey);
                String mimeType = Wailsmobile.requestAssetMime(assetKey);
                if (mimeType == null || mimeType.isEmpty()) mimeType = "text/plain";

                return new WebResourceResponse(mimeType, "UTF-8", new ByteArrayInputStream(fileBytes));
            }

            @Override
            public void onPageStarted(WebView view, String url, android.graphics.Bitmap favicon) {
                super.onPageStarted(view, url, favicon);
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

        mWebView.loadUrl("https://wails.local/");
        startEventPolling();

        // Handle initial intent for Deep Linking
        handleIntent(getIntent());
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        setIntent(intent);
        handleIntent(intent);
    }

    private void handleIntent(Intent intent) {
        if (intent == null) return;
        WailsApplication app = (WailsApplication) getApplication();
        for (WailsPlugin plugin : app.getPlugins().values()) {
            plugin.onNewIntent(intent);
        }
    }

    private void startEventPolling() {
        new Thread(() -> {
            while (mIsPolling) {
                String eventJson = Wailsmobile.pollNativeEvent();
                if (eventJson != null && !eventJson.isEmpty()) {
                    mHandler.post(() -> {
                        String script = "if(window.WailsBind && window.WailsBind.dispatch) { " +
                                "window.WailsBind.dispatch(" + eventJson + "); " +
                                "} else if(window.WailsEvents && window.WailsEvents.dispatch) { " +
                                "window.WailsEvents.dispatch(" + eventJson + "); " +
                                "}";
                        mWebView.evaluateJavascript(script, null);
                    });
                }
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    break;
                }
            }
        }).start();
    }

    @Override
    protected void onPause() {
        super.onPause();
        if (mWebView != null) mWebView.onPause();
        Wailsmobile.handleNativeAction("lifecycle:pause", "");
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (mWebView != null) mWebView.onResume();
        Wailsmobile.handleNativeAction("lifecycle:resume", "");
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
        WailsApplication app = (WailsApplication) getApplication();
        for (WailsPlugin plugin : app.getPlugins().values()) {
            plugin.onRequestPermissionsResult(requestCode, permissions, grantResults);
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        WailsApplication app = (WailsApplication) getApplication();
        for (WailsPlugin plugin : app.getPlugins().values()) {
            plugin.onActivityResult(requestCode, resultCode, data);
        }
    }

    @Override
    protected void onDestroy() {
        mIsPolling = false;
        if (mWebView != null) mWebView.destroy();
        super.onDestroy();
    }
}
