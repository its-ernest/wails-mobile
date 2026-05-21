package com.example.wailsmobile;

import android.annotation.SuppressLint;
import android.os.Bundle;
import android.util.Log;
import android.webkit.JavascriptInterface;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import androidx.appcompat.app.AppCompatActivity;

import java.io.ByteArrayInputStream;
import java.nio.charset.StandardCharsets;

import wailsmobile.WailsMobile;

public class WailsWebViewActivity extends AppCompatActivity {

    private WebView mWebView;

    @SuppressLint("SetJavaScriptEnabled")
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        String initResult = WailsMobile.startApplication();
        Log.d("WailsMobile", "Go backend init result: " + initResult);

        mWebView = new WebView(this);
        setContentView(mWebView);

        mWebView.getSettings().setJavaScriptEnabled(true);
        mWebView.getSettings().setDomStorageEnabled(true);

        mWebView.addJavascriptInterface(new Object() {
            @JavascriptInterface
            public String callGo(String methodKey, String jsonArgsPayload) {
                return WailsMobile.handleMessageFromFrontend(methodKey, jsonArgsPayload);
            }

            @JavascriptInterface
            public String callNative(String methodKey, String jsonArgsPayload) {
                return WailsMobile.handleNativeAction(methodKey, jsonArgsPayload);
            }
        }, "WailsBind");

        mWebView.setWebViewClient(new WebViewClient() {
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                String urlPath = request.getUrl().getPath();
                if (urlPath == null || urlPath.equals("/")) {
                    urlPath = "index.html";
                }

                byte[] fileBytes = WailsMobile.requestAssetBytes(urlPath);
                String mimeType = WailsMobile.requestAssetMime(urlPath);
                if (mimeType == null || mimeType.isEmpty()) {
                    mimeType = "text/plain";
                }

                return new WebResourceResponse(
                        mimeType,
                        "UTF-8",
                        new ByteArrayInputStream(fileBytes)
                );
            }
        });

        String html = new String(WailsMobile.requestAssetBytes("index.html"), StandardCharsets.UTF_8);
        mWebView.loadDataWithBaseURL("https://wails.local/", html, "text/html", "UTF-8", null);
    }
}
