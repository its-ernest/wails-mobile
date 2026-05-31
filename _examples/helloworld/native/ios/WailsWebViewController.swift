import UIKit
import WebKit
import Wailsmobile

final class WailsWebViewController: UIViewController, WKNavigationDelegate {
    private lazy var webView: WKWebView = {
        let configuration = WKWebViewConfiguration()
        configuration.setURLSchemeHandler(WailsURLSchemeHandler(), forURLScheme: "wails")

        let view = WKWebView(frame: .zero, configuration: configuration)
        view.navigationDelegate = self
        view.allowsBackForwardNavigationGestures = true
        view.configuration.preferences.javaScriptEnabled = true
        view.configuration.preferences.javaScriptCanOpenWindowsAutomatically = true
        return view
    }()

    private var eventTimer: Timer?

    override func loadView() {
        view = webView
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        loadRootAsset()
        startEventPolling()
    }

    override func viewDidAppear(_ animated: Bool) {
        super.viewDidAppear(animated)
        Wailsmobile.handleNativeAction("lifecycle:resume", "")
    }

    override func viewWillDisappear(_ animated: Bool) {
        super.viewWillDisappear(animated)
        Wailsmobile.handleNativeAction("lifecycle:pause", "")
    }

    private func loadRootAsset() {
        guard let url = URL(string: "wails://wails.local/index.html") else {
            return
        }
        webView.load(URLRequest(url: url))
    }

    private func startEventPolling() {
        eventTimer = Timer.scheduledTimer(withTimeInterval: 0.1, repeats: true) { [weak self] _ in
            guard let self = self else { return }
            let eventJson = Wailsmobile.pollNativeEvent()
            if !eventJson.isEmpty {
                self.dispatchEvent(eventJson)
            }
        }
    }

    private func dispatchEvent(_ eventJson: String) {
        let script = "if(window.WailsBind && window.WailsBind.dispatch) { window.WailsBind.dispatch(") + eventJson + "); }"
        webView.evaluateJavaScript(script, completionHandler: nil)
    }

    func webView(_ webView: WKWebView, didFinish navigation: WKNavigation!) {
        injectWailsSupportJS()
    }

    private func injectWailsSupportJS() {
        let js = "window.WailsEvents = window.WailsEvents || { listeners: {}, on: function(name, cb) { if(!this.listeners[name]) this.listeners[name] = []; this.listeners[name].push(cb); }, dispatch: function(obj) { var name = obj.name; var data = obj.data; if(this.listeners[name]) { this.listeners[name].forEach(function(cb) { try { cb(data); } catch(e) { console.error(e); } }); } } };"
        webView.evaluateJavaScript(js, completionHandler: nil)
    }
}
