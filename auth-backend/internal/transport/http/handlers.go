package httpserver

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"service/internal/service"
	pb "service/pkg/grpc/auth_v1"
)

type Server struct {
	Service *service.Service
}

func NewServer(service *service.Service) *Server {
	return &Server{
		Service: service,
	}
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на логин от %s", r.RemoteAddr)

	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:3001"
	}
	log.Printf("Redirect URI: %s", redirectURI)

	var req pb.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Ошибка при декодировании запроса логина: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Попытка логина пользователя: %s", req.Name)

	resp, err := s.Service.Login(r.Context(), &req)
	if err != nil {
		log.Printf("Ошибка логина для пользователя %s: %v", req.Name, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Успешный логин пользователя: %s", req.Name)

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		log.Printf("Неверный redirect URI: %s, ошибка: %v", redirectURI, err)
		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}

	redirectURL.Fragment = url.Values{"token": {resp.Token}}.Encode()

	log.Printf("Перенаправление пользователя %s на URL: %s", req.Name, redirectURL.String())

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

// Обработчик для регистрации
func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос на регистрацию от %s", r.RemoteAddr)

	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = "http://localhost:3001"
	}
	log.Printf("Redirect URI: %s", redirectURI)

	var req pb.RegisterRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Ошибка при декодировании запроса регистрации: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Не рекомендуется логировать чувствительную информацию, такую как пароли
	log.Printf("Попытка регистрации пользователя: %s", req.Name)

	resp, err := s.Service.Register(r.Context(), &req)
	if err != nil {
		log.Printf("Ошибка регистрации для пользователя %s: %v", req.Name, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Успешная регистрация пользователя: %s", req.Name)

	redirectURL, err := url.Parse(redirectURI)
	if err != nil {
		log.Printf("Неверный redirect URI: %s, ошибка: %v", redirectURI, err)
		http.Error(w, "Invalid redirect URI", http.StatusBadRequest)
		return
	}

	redirectURL.Fragment = url.Values{"token": {resp.Token}}.Encode()

	log.Printf("Перенаправление пользователя %s на URL: %s", req.Name, redirectURL.String())

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}
