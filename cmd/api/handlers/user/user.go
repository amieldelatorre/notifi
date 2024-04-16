package user // import "github.com/amieldelatorre/notifi/cmd/api/handlers/user"

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amieldelatorre/notifi/cmd/api/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	getUserHandler := http.HandlerFunc(getUser)

	mux.Handle("GET /api/v1/user", middleware.ApiKeyAuth(getUserHandler))
	mux.HandleFunc("POST /api/v1/user", postUser)
	mux.HandleFunc("PUT /api/v1/user", putUser)
	mux.HandleFunc("DELETE /api/v1/user", deleteUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Println(r.Header.Get("x-user-id"))

	data := map[string]string{
		"method": "GET",
	}

	json.NewEncoder(w).Encode(data)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	data := map[string]string{
		"method": "POST",
	}

	json.NewEncoder(w).Encode(data)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	data := map[string]string{
		"method": "PUT",
	}

	json.NewEncoder(w).Encode(data)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)

	data := map[string]string{
		"method": "DELETE",
	}

	json.NewEncoder(w).Encode(data)
}
