import UIKit

/// WailsPlugin is the iOS equivalent of the Android `com.wailsplugin.WailsPlugin` interface.
public protocol WailsPlugin {
    /// Unique plugin domain used by Go-to-native calls, e.g. "permissions".
    var domain: String { get }

    /// Called when the plugin is attached to the top-level view controller.
    func onAttach(_ viewController: UIViewController)

    /// Called by the native handler when Go needs the plugin to perform an action.
    func handleAction(_ action: String, jsonArgsPayload: String) -> String

    /// Optional lifecycle hook for deep links or open URL events.
    func onOpenURL(_ url: URL)

    /// Optional lifecycle hook for permission results.
    func onPermissionResult(_ requestCode: Int, permissions: [String], grantResults: [Int])
}

public extension WailsPlugin {
    func onAttach(_ viewController: UIViewController) {}
    func onOpenURL(_ url: URL) {}
    func onPermissionResult(_ requestCode: Int, permissions: [String], grantResults: [Int]) {}
}
