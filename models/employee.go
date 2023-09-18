package models

type Employee struct {
	Id             string `json:"Id"`
	Surname        string `json:"Surname"`
	Name           string `json:"Name"`
	Patronymic     string `json:"Patronymic"`
	EmployeeNumber *int   `json:"EmployeeNumber"`
}
