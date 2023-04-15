package server

func (s *Server) routes() {
	api := s.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", s.handleUserPost()).Methods("POST")

	s.Router.HandleFunc("/ws", s.handleWs())
}
