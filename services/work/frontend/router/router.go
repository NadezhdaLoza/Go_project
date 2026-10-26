package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/work/frontend/handlers"
	"github.com/rs/cors"
)

func GetRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/work/works", handlers.WorksHandler)

	router.HandleFunc("/work/works/create", handlers.CreateWorkGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/work/works/create", handlers.CreateWorkPostHandler).Methods(http.MethodPost)

	router.HandleFunc("/work/works/update/{id}", handlers.UpdateWorkGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/work/works/update", handlers.UpdateWorkPostHandler).Methods(http.MethodPost)

	router.HandleFunc("/work/works/delete", handlers.DeleteWorkHandler).Methods(http.MethodPost)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../../../static"))))
	router.PathPrefix("/work/static/").Handler(http.StripPrefix("/work/static/", http.FileServer(http.Dir("../static"))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
	})

	handler := c.Handler(router)
	return handler
}