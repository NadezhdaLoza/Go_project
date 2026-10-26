package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nadezhda/documents/services/department/frontend/handlers"
	"github.com/rs/cors"
)

func GetRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/department/departments", handlers.DepartmentsHandler)

	router.HandleFunc("/department/departments/create", handlers.CreateDepartmentGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/department/departments/create", handlers.CreateDepartmentPostHandler).Methods(http.MethodPost)

	router.HandleFunc("/department/departments/update/{id}", handlers.UpdateDepartmentGetHandler).Methods(http.MethodGet)
	router.HandleFunc("/department/departments/update", handlers.UpdateDepartmentPostHandler).Methods(http.MethodPost)

	router.HandleFunc("/department/departments/delete", handlers.DeleteDepartmentHandler).Methods(http.MethodPost)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../../../static"))))
	router.PathPrefix("/department/static/").Handler(http.StripPrefix("/department/static/", http.FileServer(http.Dir("../static"))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost},
	})

	handler := c.Handler(router)
	return handler
}