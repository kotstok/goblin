package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kotstok/goblin/internal/app/ws"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	config *Config
	router *mux.Router
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()
	s.configureWebSocket()

	log.Println(fmt.Sprintf("Starting Server http://localhost%s", s.config.BindAddr))

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureRouter() {
	//file server
	s.router.Handle("/img/{rest}", http.StripPrefix("/img/", http.FileServer(http.Dir("web/public/img/"))))
	s.router.PathPrefix("/js/{rest}").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/public/js/"))))
	s.router.PathPrefix("/css/{rest}").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/public/css/"))))

	// chat router
	s.router.HandleFunc("/", s.handleIndex())
}

func (s *Server) configureWebSocket() {

	hub := ws.Start()

	log.Println("Websocket Start: ok")

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

		outputHTML(w, "index.html", nil)
	}
}
