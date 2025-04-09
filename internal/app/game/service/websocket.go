package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/Echin-h/HangZhou-Monopoly/internal/app/game/model"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境应该检查来源
	},
}

type Client struct {
	conn   *websocket.Conn
	TeamID string
	send   chan []byte
	game   *model.GameState
}

type GameServer struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	game       *model.GameState
	mu         sync.Mutex
}

func NewGameServer(game *model.GameState) *GameServer {
	return &GameServer{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		game:       game,
	}
}

func (s *GameServer) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			log.Println("Client registered:", client.TeamID)
			//s.notifyGameState()

		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				close(client.send)
				delete(s.clients, client)
				log.Println("Client unregistered:", client.TeamID)
			}

		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(s.clients, client)
				}
			}
		}
	}
}

func (s *GameServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	teamID := r.URL.Query().Get("team_id")
	if teamID == "" {
		conn.WriteMessage(websocket.CloseMessage, []byte("team_id is required"))
		conn.Close()
		return
	}

	client := &Client{
		conn:   conn,
		TeamID: teamID,
		send:   make(chan []byte, 256),
		game:   s.game,
	}

	s.register <- client

	//go s.handleMessages(client)
	//go s.writePump(client)
}
