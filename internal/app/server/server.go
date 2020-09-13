package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kotstok/goblin/internal/app/ws"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

type Server struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()
	s.configureWebSocket()

	s.logger.Info(fmt.Sprintf("Starting Server http://localhost%s", s.config.BindAddr))

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *Server) configureRouter() {
	// ..
	s.router.HandleFunc("/", s.handleIndex())
}

func (s *Server) configureWebSocket() {

	hub := ws.Start()

	logrus.Info("Websocket Start: ok")

	s.router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRegister(hub, w, r)
	})
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles("web/templates/" + filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	// ..
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		params := map[string]interface{}{"Title": "Goblin"}

		outputHTML(w, "index.html", params)
	}
}
