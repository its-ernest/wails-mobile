import UIKit

public final class ExamplePlugin: WailsPlugin {
    public let domain = "example"

    public init() {}

    public func onAttach(_ viewController: UIViewController) {
        // Attach the plugin to the active view controller for UI or permissions work.
    }

    public func handleAction(_ action: String, jsonArgsPayload: String) -> String {
        switch action {
        case "echo":
            return "{\"result\": \"\(jsonArgsPayload)\"}"
        default:
            return "{\"error\": \"Unknown action \(action)\"}"
        }
    }
}
