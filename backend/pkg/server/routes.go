package server

func (s *Server) routes() {
	api := s.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", s.handleUserPost()).Methods("POST")
	api.HandleFunc("/room", s.handleRoomPost()).Methods("POST")
	api.HandleFunc("/room/join", s.handleRoomJoinPost()).Methods("POST")

	s.Router.HandleFunc("/ws", s.handleWs())
}
