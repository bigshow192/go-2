package service

import (
	"context"
	"database/sql"
	. "go_second/employee"
	"reflect"

	sv "github.com/core-go/core"

	//. "github.com/core-go/search"
	q "github.com/core-go/sql"
)

type employeeService struct {
	db         *sql.DB
	repository sv.ViewRepository
}
type EmployeeService interface {
	All(ctx context.Context) (*[]Employee, error)
	Create(ctx context.Context, employee Employee) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
	Update(ctx context.Context, employee Employee) (int64, error)
	Patch(ctx context.Context, employee map[string]interface{}) (int64, error)
	Load(ctx context.Context, id string) (int64, error)
	AllRepo(ctx context.Context) (interface{}, error)
}

func NewEmployeeService(repository sv.ViewRepository) EmployeeService {
	return &employeeService{repository: repository}
}
func (r *employeeService) AllRepo(ctx context.Context) (interface{}, error) {
	employees, err := r.repository.All(ctx)
	if err != nil {
		return nil, err
	}
	return &employees, nil
}
func (r *employeeService) Load(ctx context.Context, id string) (int64, error) {
	query := "Select * from employee where id=$1"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return -1, err

	}
	res, err := stmt.ExecContext(ctx, query, id)
	if err != nil {
		return -2, err

	}
	return res.RowsAffected()

}
func SetDB(DB *sql.DB) *employeeService {
	var e employeeService
	e.db = DB
	return &e
}
func (r *employeeService) Patch(ctx context.Context, employee map[string]interface{}) (int64, error) {
	employeeType := reflect.TypeOf(Employee{})
	jsonColumnMap := q.MakeJsonColumnMap(employeeType)
	colMap := q.JSONToColumns(employee, jsonColumnMap)
	keys, _ := q.FindPrimaryKeys(employeeType)
	query, args := q.BuildToPatch("employee", colMap, keys, q.BuildDollarParam)
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
func (r *employeeService) Create(ctx context.Context, employee Employee) (int64, error) {
	query := "insert into employee (id, username, email,phone) values ($1,$2,$3,$4)"
	//stmt := fmt.Sprintf("INSERT INTO employee(id, username,email,phone) VALUES('%s','%s','%s','%s')", employee.Id, employee.Username, employee.Email, employee.Phone)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return -2, err
	}
	res, er1 := stmt.ExecContext(ctx, employee.Id, employee.Username, employee.Email, employee.Phone)
	if er1 != nil {
		return -1, er1
	}
	return res.RowsAffected()

}
func (r *employeeService) Delete(ctx context.Context, id string) (int64, error) {
	//stmt := fmt.Sprintf("Delete from employee where id='%s'", id)
	query := "Delete from employee where id=$1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
func (r *employeeService) Update(ctx context.Context, employee Employee) (int64, error) {
	//stmt := fmt.Sprintf("Update employee set username='%s', email='%s', phone='%s' where id='%s';", employee.Username, employee.Email, employee.Phone, employee.Id)
	//fmt.Println(stmt)
	query := "Update employee set username=$1, email=$2, phone=$3 where id=%s;"
	res, err := r.db.ExecContext(ctx, query, employee.Username, employee.Email, employee.Phone, employee.Id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
func (r *employeeService) All(ctx context.Context) (*[]Employee, error) {
	query := "Select * from employee;"
	res, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	//fmt.Println(query)
	var employees []Employee
	for res.Next() {
		var employee Employee
		err := res.Scan(&employee.Id, &employee.Username, &employee.Email, &employee.Phone)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)

	}
	return &employees, nil

}
