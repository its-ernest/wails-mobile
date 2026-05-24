/**
 * app.js - User Application Logic
 */
function updateOutput(text) {
    document.getElementById('output').textContent = text;
}

window.addEventListener('DOMContentLoaded', () => {

    document.getElementById('send').addEventListener('click', async () => {
        const name = document.getElementById('name').value;
        try {
            const result = await Wails.CallGo('AppService.SayHello', name);
            updateOutput(result.message);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('check-perm').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('PermissionPlugin.Check', "android.permission.CAMERA");
            // Go returns the raw string from Java
            const parsed = JSON.parse(result);
            updateOutput(`Camera Status: ${parsed.status}`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('check-notif').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('PermissionPlugin.Check', "android.permission.POST_NOTIFICATIONS");
            // Go returns the raw string from Java
            const parsed = JSON.parse(result);
            updateOutput(`Notif Status: ${parsed.status}`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('request-perm').addEventListener('click', async () => {
        try {
            await Wails.CallGo('PermissionPlugin.Request', "android.permission.CAMERA");
            updateOutput("Requested camera... check the system dialog.");
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('request-notif').addEventListener('click', async () => {
        try {
            await Wails.CallGo('PermissionPlugin.Request', "android.permission.POST_NOTIFICATIONS");
            updateOutput("Requested notification... check the system dialog.");
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('ping').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('AppService.Ping');
            updateOutput(`Bridge Ping: ${result.status}`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('enqueue-work').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('AppService.EnqueuePeriodic');
            const parsed = JSON.parse(result);
            updateOutput(`Work Enqueued: ${parsed.status} (ID: ${parsed.id})`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    document.getElementById('notify').addEventListener('click', async () => {
        try {
            // Passing ID: 0 (or omitting it) will now generate a unique notification instead of overwriting
            const result = await Wails.CallGo('NotificationPlugin.Post', {
                id: 30,
                title: "Unique Notification",
                body: "This won't overwrite previous ones unless they share the same ID.",
                importance: "HIGH"
            });
            const parsed = JSON.parse(result);
            updateOutput(`Notification Posted: ID ${parsed.id}`);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    // 4. Event Listener
    Wails.on('permissions:changed', (data) => {
        console.log("EVENT RECEIVED:", data);
        updateOutput(`ASYNC EVENT: Permission for ${data.permission} was ${data.granted ? 'GRANTED' : 'DENIED'}`);
    });
});
