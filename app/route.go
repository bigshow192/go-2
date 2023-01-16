package app

import (
	"database/sql"
	"encoding/json"

	//. "go_second/employee"
	"go_second/handler"
	mi "go_second/middleware"
	"log"
	"net/http"

	//"go_second/service"

	"github.com/gorilla/mux"
)

type Router struct {
	employeeHandler handler.EmployeeHandler
}

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
func (r *Router) Router(db *sql.DB, e handler.EmployeeHandler) {
	//var e handler.EmployeeHandler
	// var s service.EmployeeService
	// s = service.SetDB(db)
	//e = handler.SetEmployeeService(s)
	//e=handler.EmployeeHandler.
	m := mux.NewRouter()
	employee := "/employee"
	m.HandleFunc(employee, e.AllRepo).Methods("GET")
	m.HandleFunc("/", home).Methods("GET")
	m.Use(mi.LoggingMiddleware)
	m.HandleFunc(employee, e.Create).Methods("POST")
	m.HandleFunc(employee+"/{id}", e.Put).Methods("PUT")
	m.HandleFunc(employee+"/{id}", e.Delete).Methods("Delete")
	m.HandleFunc(employee+"/{id}", e.Patch).Methods("PATCH")
	m.HandleFunc(employee+"/search", e.Search).Methods("GET", "POST")
	log.Fatal(http.ListenAndServe(":5000", m))

}
