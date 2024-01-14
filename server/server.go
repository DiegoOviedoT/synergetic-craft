package server

import (
	"github.com/labstack/echo/v4"
)

type server struct {
	echo *echo.Echo
	port string
}

func New(port string) *server {
	return &server{
		echo: echo.New(),
		port: port,
	}
}

func (s *server) Start() error {
	return s.echo.Start(s.port)
}
