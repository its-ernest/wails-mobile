/**
 * app.js - Plugin Verification Logic
 */
function updateOutput(text) {
    const output = document.getElementById('output');
    output.textContent = `[${new Date().toLocaleTimeString()}] ${text}\n` + output.textContent;
}

window.addEventListener('DOMContentLoaded', () => {

    // --- Section 1: Notifications ---

    document.getElementById('request-notif').addEventListener('click', async () => {
        try {
            await Wails.CallGo('PermissionPlugin.Request', "android.permission.POST_NOTIFICATIONS");
            updateOutput("Requested notification permission.");
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('send-notif').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('NotificationPlugin.Post', {
                id: 0,
                title: "Wails Mobile",
                body: "This is a unique verification notification.",
                importance: "HIGH"
            });
            // Result is a JSON string from Go because it calls CallNativePlatform
            const parsed = typeof result === 'string' ? JSON.parse(result) : result;
            updateOutput(`Notification posted (ID: ${parsed.id})`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Section 2: Biometrics ---

    document.getElementById('check-bio').addEventListener('click', async () => {
        try {
            const res = await Wails.CallGo('BiometricPlugin.CanAuthenticate');
            updateOutput(`Biometric Status: ${res.status} (Available: ${res.can_authenticate})`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('auth-bio').addEventListener('click', async () => {
        try {
            updateOutput("Starting biometric authentication...");
            await Wails.CallGo('BiometricPlugin.Authenticate', {
                title: "Verify Identity",
                subtitle: "Confirm biometric to proceed",
                description: "This test ensures the native prompt bridge is working.",
                negative_button_text: "Cancel"
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Section 3: File Picker ---

    document.getElementById('pick-image').addEventListener('click', async () => {
        try {
            updateOutput("Opening image picker...");
            await Wails.CallGo('FilePickerPlugin.PickFile', {
                mime_type: "image/*",
                multiple: true
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('pick-files').addEventListener('click', async () => {
        try {
            updateOutput("Opening multi-file picker...");
            await Wails.CallGo('FilePickerPlugin.PickFile', {
                mime_type: "*/*",
                multiple: true
            });
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });


    // --- Event Listeners ---

    Wails.on('biometric:result', (result) => {
        updateOutput(`BIOMETRIC EVENT: ${result.status} (Success: ${result.success})`);
        if (result.error) updateOutput(`Detail: ${result.error}`);
    });

    Wails.on('filepicker:result', (result) => {
        if (result.error) {
            updateOutput(`FILEPICKER EVENT: ${result.error}`);
            return;
        }
        if (result.multiple) {
            updateOutput(`FILEPICKER EVENT: Selected ${result.uris.length} files`);
        } else {
            updateOutput(`FILEPICKER EVENT: Selected ${result.uri}`);
        }
    });

    Wails.on('permissions:changed', (data) => {
        updateOutput(`PERMISSION EVENT: ${data.permission} is ${data.granted ? 'GRANTED' : 'DENIED'}`);
    });

});
