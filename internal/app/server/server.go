package server

import (
	"encoding/json"
	"fmt"
	"github.com/adam-hanna/jwt-auth/jwt"
	"github.com/gorilla/mux"
	"github.com/kotstok/goblin/internal/app/errors"
	"github.com/kotstok/goblin/internal/app/ws"
	"html/template"
	"log"
	"net/http"
	"time"
)

var restrictedRoute jwt.Auth

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

	authErr := jwt.New(&restrictedRoute, jwt.Options{
		SigningMethodString:   "RS256",
		PrivateKeyLocation:    "configs/keys/app.rsa",     // `$ openssl genrsa -out app.rsa 2048`
		PublicKeyLocation:     "configs/keys/app.rsa.pub", // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
		RefreshTokenValidTime: 72 * time.Hour,
		AuthTokenValidTime:    15 * time.Minute,
		Debug:                 true,
		IsDevEnv:              true,
	})

	if authErr != nil {
		log.Println("Error initializing the JWT's!")
		log.Fatal(authErr)
	}

	s.configureRouter()
	s.configureWebSocket()

	log.Println(fmt.Sprintf("Starting Server http://localhost%s", s.config.BindAddr))

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureRouter() {
	// Error handlers  [ http errors ]
	s.router.MethodNotAllowedHandler = errors.MethodNotAllowedHandler()
	s.router.NotFoundHandler = errors.NotFoundHandler()

	//file server
	s.router.Handle("/img/{n}", http.StripPrefix("/img/", http.FileServer(http.Dir("web/public/img/"))))
	s.router.PathPrefix("/js/{n}").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/public/js/"))))
	s.router.PathPrefix("/css/{n}").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/public/css/"))))

	// chat router
	s.router.Handle("/", restrictedRoute.Handler(s.handleIndex())).Methods("GET")
	s.router.HandleFunc("/login", s.handleLogin()).Methods("GET")
	s.router.HandleFunc("/auth", s.handleJsonAuth()).Methods("POST")
}

func (s *Server) configureWebSocket() {

	hub := ws.Start()

	log.Println("Websocket Start: ok")

	s.router.Handle("/ws", restrictedRoute.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRegister(hub, w, r)
	})))
}

func outputHTML(w http.ResponseWriter, filename string, data interface{}) {
	tRoot := "web/templates/"
	t, err := template.ParseFiles(tRoot+filename, tRoot+"base.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	// .. INDEX page
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		outputHTML(w, "index.html", nil)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	// .. LOGIN page
	return func(w http.ResponseWriter, r *http.Request) {
		outputHTML(w, "login.html", nil)
	}
}

func (s *Server) handleJsonAuth() http.HandlerFunc {
	// .. JSON auth page

	type Profile struct {
		Name    string
		Hobbies []string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		profile := Profile{"Alex", []string{"snowboarding", "programming"}}

		js, err := json.Marshal(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
