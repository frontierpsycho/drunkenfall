package websockets

import (
	"log"

	"golang.org/x/net/websocket"
	"time"
)

// Server represents Websocket server
type Server struct {
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// NewServer creates new chat server
func NewServer() *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}

// Add adds a new client
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Del deletes a client
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// SendAll sends a broadcast message
func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

// Done closes the server
func (s *Server) Done() {
	s.doneCh <- true
}

// Err sends an error
func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	// We only need to send the last message, and only if something actually
	// exists.
	if len(s.messages) != 0 {
		c.Write(s.messages[len(s.messages)-1])
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		go c.Write(msg)
	}
}

// Pulse sends a pulse message to keep the connections alive
func (s *Server) Pulse() {
	for _, c := range s.clients {
		go c.Ping()
	}
}

// OnConnected is the function to be passed to http.Handle(), wrapped in a
// websocket.Handler().
func (s *Server) OnConnected(ws *websocket.Conn) {
	defer func() {
		log.Print("Deferring")
		err := ws.Close()
		if err != nil {
			s.errCh <- err
		}
	}()

	client := NewClient(ws, s)
	s.Add(client)
	client.Listen()
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Websocket handler initialized")

	// Setup the worst Ping implentation of all time
	go func() {
		for {
			s.Pulse()
			time.Sleep(time.Second * 55)
		}
	}()

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
