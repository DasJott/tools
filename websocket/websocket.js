// WebSocket
// simply include this file and call from your code:
// setupWebSocket(uri)
// where uri is your handshake endpoint.
// example:
// uri = 'wss://' + window.location.host + "/ws";

function setupWebSocket(uri) {
    if (!"WebSocket" in window) {
        output("Use a better browser!");
        return;
    }

    var socket = {
        register: function() {
            var ws = new WebSocket(uri);
            ws.onopen = function() {
                console.log("socket connection opened.");
            };
            ws.onmessage = function (evt) {
                console.log(evt.data);
            };
            ws.onclose = function() {
                console.log("socket connection closed. reconnecting...");
                setTimeout(function() {
                    socket.register();
                }, 1000);
            };
        }
    }

    socket.register();
}
