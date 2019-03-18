// date: 2019-03-15
package quote

import "net/http"

type Server struct {
	ServerN
}

func (s *Server) Listen() {
	http.Handle("/", s.newServerHandler())
	http.ListenAndServe(":12345", nil)
}

func New(s ServerN) *Server {
	return &Server{s}
}
