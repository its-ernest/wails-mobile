async function CallGoBackend(methodName, ...args) {
    if (!window.WailsBind) {
        throw new Error('Native framework container bridge not found.');
    }

    const jsonArgs = JSON.stringify(args);
    const rawResponse = window.WailsBind.callGo(methodName, jsonArgs);
    const parsed = JSON.parse(rawResponse);

    if (parsed.error) {
        throw new Error(parsed.error);
    }

    return parsed.result;
}

async function loadHelloWorld() {
    const nameInput = document.getElementById('name');
    const output = document.getElementById('output');

    output.textContent = 'Calling Go backend...';

    try {
        const result = await CallGoBackend('HelloService.SayHello', nameInput.value.trim());
        output.textContent = result.message;
    } catch (err) {
        output.textContent = `Error: ${err.message}`;
        console.error('Bridge error:', err);
    }
}

window.addEventListener('DOMContentLoaded', () => {
    document.getElementById('send').addEventListener('click', loadHelloWorld);
});
