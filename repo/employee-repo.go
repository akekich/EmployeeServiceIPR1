package repo

import (
	_ "embed"
	"employee-service/models"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type EmployeeRepo interface {
	AddEmployee(employee models.Employee)
	GetEmployees() []models.Employee
}

type SqlRepository struct {
	db *sqlx.DB
}

func NewSqlRepository(connectionString string) SqlRepository {
	conf, err := pgx.ParseConnectionString(connectionString)
	if err != nil {
		fmt.Println("Ошибка конфига")
		return SqlRepository{nil}
	}

	conn := sqlx.NewDb(stdlib.OpenDB(conf), "pgx")

	return SqlRepository{db: conn}
}

//go:embed sql/add_employee.sql
var addEmployee string

func (s SqlRepository) AddEmployee(employee models.Employee) {
	_, err := s.db.Query(addEmployee, employee.Surname, employee.Name, employee.Patronymic, employee.EmployeeNumber)

	if err != nil {
		fmt.Println(fmt.Sprintf("Ошибка при записи сотрудника %s", err))
	}
}

//go:embed sql/get_employee.sql
var getEmployee string

func (s SqlRepository) GetEmployees() []models.Employee {
	var employees []models.Employee
	rows, err := s.db.Query(getEmployee)

	if err != nil {
		fmt.Println(fmt.Sprintf("Ошибка при записи сотрудника %s", err))
	}

	for rows.Next() {
		employee := models.Employee{}
		var (
			Surname        string
			Name           string
			Patronymic     string
			EmployeeNumber *int
		)

		if err = rows.Scan(&Surname, &Name, &Patronymic, &EmployeeNumber); err != nil {
			continue
		}

		employee.Surname = Surname
		employee.Name = Name
		employee.Patronymic = Patronymic
		employee.EmployeeNumber = EmployeeNumber

		employees = append(employees, employee)
	}

	return employees
}
