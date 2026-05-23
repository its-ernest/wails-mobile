/**
 * app.js - User Application Logic
 */
function updateOutput(text) {
    document.getElementById('output').textContent = text;
}

window.addEventListener('DOMContentLoaded', () => {
    // 1. Standard Go Call
    document.getElementById('send').addEventListener('click', async () => {
        const name = document.getElementById('name').value;
        try {
            const result = await Wails.CallGo('AppService.SayHello', name);
            updateOutput(result.message);
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    // 2. Permission Plugin (Managed by Go)
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

    document.getElementById('request-perm').addEventListener('click', async () => {
        try {
            await Wails.CallGo('PermissionPlugin.Request', "android.permission.CAMERA");
            updateOutput("Requested camera... check the system dialog.");
        } catch (err) {
            updateOutput(`Error: ${err.message}`);
        }
    });

    // 3. Native Go Method
    document.getElementById('ping').addEventListener('click', async () => {
        try {
            const result = await Wails.CallGo('AppService.Ping');
            updateOutput(`Bridge Ping: ${result.status}`);
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
