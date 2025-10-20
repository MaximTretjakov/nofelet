package httpserver

import "time"

type Option func(s *Server)

func WithReadTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.server.ReadTimeout = d
	}
}

func WithReadHeaderTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.server.ReadHeaderTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.server.WriteTimeout = d
	}
}

func WithShutdownTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = d
	}
}

func WithAddress(address string) Option {
	return func(s *Server) {
		s.server.Addr = address
	}
}
