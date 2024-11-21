package main

import (
	"github.com/rs/cors"
)

//func allowCORS(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
//		w.Header().Set("Access-Control-Allow-Credentials", "true")
//
//		if r.Method == http.MethodOptions {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		h.ServeHTTP(w, r)
//	})
//}

func allowCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // ваш фронтенд адрес
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}
