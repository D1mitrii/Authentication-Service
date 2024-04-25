package httpserver

import (
	"fmt"
	"net"
	"time"
)

type Option func(*Server)

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

func Port(port int) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", fmt.Sprint(port))
	}
}
