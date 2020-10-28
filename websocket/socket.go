package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// MessageType is the type of message to be sent through the websocket
type MessageType int

const (
	// Text represents a text message
	Text = websocket.TextMessage
	// Binary represents a binary message
	Binary = websocket.BinaryMessage
)

// Socket is a simple and easy to use websocket implementation.
// for a new instance use NewSocket
type Socket struct {
	Receive     func(MessageType, []byte)
	CheckOrigin func(r *http.Request) bool

	// sc         chan []byte
	upgrader   *websocket.Upgrader
	connection *websocket.Conn
}

// New returns a pointer to a new Socket object.
// read and write are the buffer sizes of each.
func New(read, write int) *Socket {
	s := &Socket{}
	s.upgrader = &websocket.Upgrader{
		ReadBufferSize:  read,
		WriteBufferSize: write,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	// s.sc = make(chan []byte)

	return s
}

func (s *Socket) read() {
	for {
		dataType, data, err := s.connection.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		s.Receive(MessageType(dataType), data)
	}
}

// Send can be used to send data to the connection
func (s *Socket) Send(dataType MessageType, data []byte) {
	if err := s.connection.WriteMessage(int(dataType), data); err != nil {
		fmt.Println(err.Error())
	}
}

// Handle is used within a HTTP handler
func (s *Socket) Handle(w http.ResponseWriter, r *http.Request) (err error) {
	s.Close()

	if s.CheckOrigin != nil {
		s.upgrader.CheckOrigin = s.CheckOrigin
	}

	s.connection, err = s.upgrader.Upgrade(w, r, nil)
	if err == nil && s.Receive != nil {
		go s.read()
	}

	return err
}

// Close closes the connection of the websocket
func (s *Socket) Close() error {
	if s.connection != nil {
		return s.connection.Close()
	}
	return nil
}
