package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) handleWs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("addr", r.RemoteAddr)
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error(err)
			w.Write([]byte(err.Error()))
			return
		}
		log.Info("Client connected")
		go reader(ws)
	}
}

func reader(conn *websocket.Conn) {
	defer log.Info("Client disconnected")
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithField("error", err).Warn("error reading message from ws")
			}
			return
		}

		log.WithFields(log.Fields{
			"addr":        conn.RemoteAddr(),
			"messageType": messageType,
			"message":     string(p),
		}).Info("Client Message")

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Warn(err)
			return
		}
	}
}
