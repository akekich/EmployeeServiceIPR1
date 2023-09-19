package handlers

import (
	"employee-service/models"
	"employee-service/repo"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type EmployeeHandler interface {
	AddHandler(w http.ResponseWriter, r *http.Request)
	GetHandler(w http.ResponseWriter, r *http.Request)
}

type HandlerEmployee struct {
	repo repo.EmployeeRepo
}

func NewHandler(employees repo.EmployeeRepo) HandlerEmployee {
	return HandlerEmployee{repo: employees}
}

func (e HandlerEmployee) AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Ошибка чтения тела запроса")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		fmt.Println(fmt.Sprintf("Ошибка %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	e.repo.AddEmployee(employee)
}

func (e HandlerEmployee) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	employees := e.repo.GetEmployees()
	if len(employees) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	vars := mux.Vars(r)
	sortName := vars["sort"]
	sortedEmployees := []models.Employee{}
	employeesWithNumber := []models.Employee{}
	employeesWithoutNumber := []models.Employee{}

	for _, employee := range employees {
		if employee.EmployeeNumber != nil {
			employeesWithNumber = append(employeesWithNumber, employee)
		} else {
			employeesWithoutNumber = append(employeesWithoutNumber, employee)
		}
	}

	switch sortName {
	case "buble":
		length := len(employeesWithNumber)
		fmt.Println(length)
		for i := 0; i < (length - 1); i++ {
			for j := 0; j < ((length - 1) - i); j++ {
				if *employeesWithNumber[j].EmployeeNumber > *employeesWithNumber[j+1].EmployeeNumber {
					employeesWithNumber[j], employeesWithNumber[j+1] = employeesWithNumber[j+1], employeesWithNumber[j]
				}
			}
		}

		sortedEmployees = append(sortedEmployees, employeesWithoutNumber...)
		sortedEmployees = append(sortedEmployees, employeesWithNumber...)
	case "quick":
		low := 0
		high := len(employeesWithNumber) - 1
		quickSort(employeesWithNumber, low, high)
		sortedEmployees = append(employeesWithoutNumber, employeesWithNumber...)
	}

	jsonData, err := json.Marshal(sortedEmployees)
	if err != nil {
		fmt.Println("Ошибка заджисонивания")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonData)
}

func quickSort(arr []models.Employee, low int, high int) {
	if low < high {
		pi := partion(arr, low, high)

		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func partion(arr []models.Employee, low int, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if *arr[j].EmployeeNumber < *pivot.EmployeeNumber {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
