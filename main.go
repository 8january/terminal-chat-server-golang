package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Room    string `json:"room"`
	Content string `json:"content"`
}

type Server struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*Room),
	}
}

type Room struct {
	connections []*Connection
	mu          sync.RWMutex
}

func NewRoom() *Room {
	return &Room{
		connections: make([]*Connection, 0),
	}
}

type Connection struct {
	conn *websocket.Conn
	send chan []byte
}

func (s *Server) getOrCreateRoom(roomName string) *Room {
	s.mu.Lock()
	defer s.mu.Unlock()

	if room, exists := s.rooms[roomName]; exists {
		return room
	}

	room := NewRoom()
	s.rooms[roomName] = room
	return room
}

func (s *Server) removeConnection(roomName string, conn *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, exists := s.rooms[roomName]; exists {
		room.removeConnection(conn)
		if len(room.connections) == 0 {
			delete(s.rooms, roomName)
		}
	}
}

func (r *Room) addConnection(conn *Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connections = append(r.connections, conn)
}

func (r *Room) removeConnection(conn *Connection) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, c := range r.connections {
		if c == conn {
			r.connections = append(r.connections[:i], r.connections[i+1:]...)
			break
		}
	}
}

func (r *Room) sendJoinMessage(conn *Connection, message Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	msg := Message{
		ID:      "0",
		Type:    "server_message",
		Name:    "server",
		Room:    message.Room,
		Content: message.Name + " joined",
	}

	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshalling error:", err)
	}
	r.broadcast(conn, jsonMessage)

}

func (r *Room) broadcast(sender *Connection, message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, conn := range r.connections {
		if conn != sender {
			select {
			case conn.send <- message:
			default:
				close(conn.send)
				r.removeConnection(conn)
			}
		}
	}
}

func readInitialMessage(conn *websocket.Conn) (Message, error) {
	var message Message
	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Println("Read error:", err)
		return message, err
	}
	err = json.Unmarshal(data, &message)
	if err != nil {
		log.Println("JSON Unmarshalling error", err)
	}
	return message, err

}

func (s *Server) handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	connection := &Connection{
		conn: conn,
		send: make(chan []byte, 256),
	}

	message, err := readInitialMessage(connection.conn)
	if err != nil {
		log.Println("Initial Message read error:", err)
		return
	}

	room := s.getOrCreateRoom(message.Room)
	room.addConnection(connection)
	room.sendJoinMessage(connection, message)

	go connection.writePump()
	connection.readPump(s, room, message.Room)

}

func (c *Connection) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *Connection) readPump(s *Server, r *Room, roomName string) {
	defer func() {
		s.removeConnection(roomName, c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("JSON unmarshalling error:", err)
			continue
		}

		log.Printf(
			"[RECEIVED FROM CLIENT]\nName: \"%s\" - Room: \"%s\" - Message: \"%s\" ",
			msg.Name,
			msg.Room,
			msg.Content,
		)
		r.broadcast(c, message)

	}

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var server = NewServer()
var mutex sync.Mutex

func main() {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := flag.String("addr", "0.0.0.0:"+port, "http service address")

	flag.Parse()
	log.SetFlags(0)

	server := NewServer()

	http.HandleFunc("/", server.handleConnection)
	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
