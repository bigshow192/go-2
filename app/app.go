package app

import (
	//"context"
	"database/sql"
	"fmt"
	"reflect"

	"go_second/service"
	. "go_second/service"

	"github.com/core-go/log"

	//"os"

	//"github.com/jackc/pgx"
	. "go_second/employee"
	. "go_second/handler"

	"github.com/core-go/search/convert"
	sq "github.com/core-go/sql"
	q "github.com/core-go/sql/query"
	"github.com/core-go/sql/template"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	CreateTable = `
create table if not exists employee (
  id varchar(40) not null,
  username varchar(120),
  email varchar(120),
  phone varchar(45),
  primary key (id)
)`
)

func NewApp(config Config) (*sql.DB, EmployeeHandler) {
	err := viper.Unmarshal(&config)
	fmt.Println(config)
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		//log.Fatal("cannot connect to db", err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	_, err = db.Exec(CreateTable)
	if err != nil {
		panic(err)
	}
	templates, err := template.LoadTemplates(template.Trim, "configs/query.xml")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var es EmployeeService
	es = SetDB(db)
	employeeType := reflect.TypeOf(Employee{})
	buildParam := q.GetBuild(db)
	employeeQuery, err := template.UseQueryWithArray(true, nil, "employee", templates, &employeeType, convert.ToMap, buildParam, pq.Array)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	employeeSearchBuilder, err := sq.NewSearchBuilderWithArray(db, employeeType, employeeQuery, pq.Array)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	getEmployee, err := sq.UseGetWithArray(db, "employees", employeeType, pq.Array)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	employeeHandler := NewEmployeeHandler(es, employeeSearchBuilder.Search, getEmployee, log.LogError, nil)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//cinemaType := reflect.TypeOf(cinema.Cinema{})
	employeeRepository, err := sq.NewRepositoryWithArray(db, "employee", employeeType, pq.Array)
	if err != nil {
		return nil, nil
	}
	employeeService := service.NewEmployeeService(employeeRepository)
	employeeHandler = NewEmployeeHandler(employeeService, employeeSearchBuilder.Search, getEmployee, log.LogError, nil)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return db, employeeHandler

}
