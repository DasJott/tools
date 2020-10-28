# WebSocket
Implementation for usage of WebSockets.<br>
Simple and quick for bidirectional communication with javascript code.<br>
Automaticly reconnecting on connection loss.

## Usage
What happens in this example?
1. We create an object that is able to receive and send messages - exactly what we want!
2. We attach a function to be called on incoming messages.
3. JavaScript is requesting the connection via GET, so we need to create the object within a GET handler.
4. We call `Handle()` with instances of `http.Response` and `http.Request`. It takes care of all the underlying WebSocket stuff.
5. Use `Send()` on the created object wherever and whenever you want.

### Go implementation
```go

// make socket accessible everywhere
var s *Socket

// example using Echo, but you can use whatever handler you want to use
func main() {
    e := echo.New()

    // create the socket object, specifying buffer sizes for read/write
    s = socket.New(1024, 1024)
    s.Receive = onMessage // incoming messages shall go into that function

    // create an endpoint for the socket handshake
    e.GET("/ws", func(c echo.Context) error {
        // Handle does everything for us
        return socket.Handle(c.Response(), c.Request())
    })

    // der Letzte macht's Licht aus
    defer s.Close()
}

// onMessage was attached above and receives messages
func onMessage(t socket.MessageType, data []byte) {
    fmt.Println(data)
}

// sendMessage sends messages to js
func sendMessage(str string) {
    if s != nil {
		s.Send(socket.Text, []byte(str))
    }
}
```

### JavaScript is this
If the server is restarted (while developing?) this code reconnects automaticly.<br>
There is a [javascript file](websocket.js) included that you can use.

```js
function setupWebSocket() {
    if (!"WebSocket" in window) {
        output("Use a better browser!");
        return;
    }

    // generate url to call
    var uri = 'wss://' + window.location.host + "/ws";

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
```