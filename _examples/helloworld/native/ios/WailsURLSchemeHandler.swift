import Foundation
import WebKit
import Wailsmobile

final class WailsURLSchemeHandler: NSObject, WKURLSchemeHandler {
    func webView(_ webView: WKWebView, start urlSchemeTask: WKURLSchemeTask) {
        guard let url = urlSchemeTask.request.url else {
            urlSchemeTask.didFailWithError(NSError(domain: "WailsScheme", code: 400, userInfo: nil))
            return
        }

        let assetPath = assetPath(for: url)
        let data = Wailsmobile.requestAssetBytes(assetPath)
        let mimeType = Wailsmobile.requestAssetMime(assetPath)
        let headers = ["Content-Type": mimeType]

        guard let response = HTTPURLResponse(url: url, statusCode: 200, httpVersion: "HTTP/1.1", headerFields: headers) else {
            urlSchemeTask.didFailWithError(NSError(domain: "WailsScheme", code: 500, userInfo: nil))
            return
        }

        urlSchemeTask.didReceive(response)
        urlSchemeTask.didReceive(data)
        urlSchemeTask.didFinish()
    }

    func webView(_ webView: WKWebView, stop urlSchemeTask: WKURLSchemeTask) {
        // No cleanup needed for static asset delivery.
    }

    private func assetPath(for url: URL) -> String {
        var path = url.path
        if path.isEmpty || path == "/" {
            path = "/index.html"
        }
        if path.hasPrefix("/") {
            path.removeFirst()
        }
        return path
    }
}
