package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func CheckOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     CheckOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := newWSClient(conn)
	s.clientsMu.Lock()
	s.clients[client] = struct{}{}
	s.clientsMu.Unlock()

	defer s.removeClient(client)

	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	go s.writePump(client)

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			return
		}
		s.handleClientMessage(client, payload)
	}
}

func (s *Server) writePump(client *wsClient) {
	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case <-client.done:
			return
		case msg := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				s.removeClient(client)
				return
			}
		case <-pingTicker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.removeClient(client)
				return
			}
		}
	}
}

func (s *Server) handleClientMessage(client *wsClient, payload []byte) {
	var message struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	}
	if err := json.Unmarshal(payload, &message); err != nil {
		return
	}

	switch message.Type {
	case "logs_subscribe":
		s.subscribeLogs(client, message.ID)
	case "logs_unsubscribe":
		s.unsubscribeLogs(client, message.ID)
	}
}

func (s *Server) subscribeLogs(client *wsClient, tunnelID string) {
	if tunnelID == "" {
		return
	}

	client.subsMu.Lock()
	if _, ok := client.logSubscriptions[tunnelID]; ok {
		client.subsMu.Unlock()
		return
	}

	logCh, unsubscribe := s.manager.SubscribeLogs(tunnelID)
	if logCh == nil {
		client.subsMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	client.logSubscriptions[tunnelID] = func() {
		cancel()
		unsubscribe()
	}
	client.subsMu.Unlock()

	go func() {
		defer s.unsubscribeLogs(client, tunnelID)
		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-logCh:
				if !ok {
					return
				}
				payload, err := MarshalJSON(map[string]interface{}{
					"type": "log",
					"id":   tunnelID,
					"line": line,
				})
				if err != nil {
					continue
				}
				if !client.enqueue(payload) {
					s.removeClient(client)
					return
				}
			}
		}
	}()
}

func (s *Server) unsubscribeLogs(client *wsClient, tunnelID string) {
	if tunnelID == "" {
		return
	}

	client.subsMu.Lock()
	cancel, ok := client.logSubscriptions[tunnelID]
	if ok {
		delete(client.logSubscriptions, tunnelID)
	}
	client.subsMu.Unlock()

	if ok {
		cancel()
	}
}

func (s *Server) removeClient(client *wsClient) {
	client.close()

	s.clientsMu.Lock()
	delete(s.clients, client)
	s.clientsMu.Unlock()

	client.subsMu.Lock()
	cancels := make([]context.CancelFunc, 0, len(client.logSubscriptions))
	for tunnelID, cancel := range client.logSubscriptions {
		cancels = append(cancels, cancel)
		delete(client.logSubscriptions, tunnelID)
	}
	client.subsMu.Unlock()

	for _, cancel := range cancels {
		cancel()
	}

	client.conn.Close()
}

func (s *Server) broadcast(msg string) {
	payload := []byte(msg)

	s.clientsMu.RLock()
	clients := make([]*wsClient, 0, len(s.clients))
	for client := range s.clients {
		clients = append(clients, client)
	}
	s.clientsMu.RUnlock()

	for _, client := range clients {
		if !client.enqueue(payload) {
			s.removeClient(client)
		}
	}
}
