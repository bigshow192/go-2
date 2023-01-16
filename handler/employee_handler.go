package handler

import (
	"context"
	"encoding/json"
	. "go_second/employee"
	. "go_second/service"
	"net/http"
	"reflect"

	//rp "github.com/core-go/core"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	"github.com/gorilla/mux"
	//u "rest-api-golang/src/utils"
)

type employeeHandler struct {
	employeeService EmployeeService
	Error           func(context.Context, string, ...map[string]interface{})
	Log             func(context.Context, string, string, bool, string) error
	*search.SearchHandler
}
type EmployeeHandler interface {
	//SetEmployeeService(es EmployeeService) employeeHandler
	Create(w http.ResponseWriter, r *http.Request)
	All(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Load(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	AllRepo(w http.ResponseWriter, r *http.Request)
}

func (e *employeeHandler) AllRepo(w http.ResponseWriter, r *http.Request) {
	employees, err := e.employeeService.AllRepo(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)

}
func NewEmployeeHandler(es EmployeeService, find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), load func(ctx context.Context, id interface{}, result interface{}) (bool, error), logError func(context.Context, string, ...map[string]interface{}), writeLog func(context.Context, string, string, bool, string) error) EmployeeHandler {
	searchModelType := reflect.TypeOf(EmployeeFilter{})
	modelType := reflect.TypeOf(Employee{})
	searchHandler := search.NewSearchHandler(find, modelType, searchModelType, logError, writeLog)
	return &employeeHandler{employeeService: es, SearchHandler: searchHandler, Error: logError, Log: writeLog}
}
func (e *employeeHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	res, err := e.employeeService.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}
func (e *employeeHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	var employee Employee
	employeeType := reflect.TypeOf(employee)
	_, jsonMap, _ := sv.BuildMapField(employeeType)
	body, er1 := sv.BuildMapAndStruct(r, &employee)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	if len(employee.Id) == 0 {
		employee.Id = id
	} else if id != employee.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	json, er2 := sv.BodyToJsonMap(r, employee, body, []string{"id"}, jsonMap)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}

	res, er3 := e.employeeService.Patch(r.Context(), json)
	if er3 != nil {
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)

}
func (e *employeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var em Employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&em); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	res, er1 := e.employeeService.Create(r.Context(), em)
	if er1 != nil {
		http.Error(w, er1.Error(), 409)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
func (e *employeeHandler) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var em Employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&em); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if id != em.Id {
		http.Error(w, "Id is wrong", 400)
		return
	}
	res, err := e.employeeService.Update(r.Context(), em)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if res == 0 {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
func (e *employeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	//var id string
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", 400)
		return
	}
	res, err := e.employeeService.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}
func (e *employeeHandler) All(w http.ResponseWriter, r *http.Request) {
	employees, err := e.employeeService.All(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)

}
func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}
