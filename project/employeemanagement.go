package project

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" 
	"github.com/gorilla/mux"           
	"github.com/joho/godotenv"         
	

	_ "Project/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)


var db *sql.DB


type Department struct {
	ID   int    `json:"id"`   
	Name string `json:"name"` 
}
type Employee struct {
	ID           int       `json:"id"`            
	Name         string    `json:"name"`          
	Email        string    `json:"email"`        
	Phone        string    `json:"phone"`         
	Salary       float64   `json:"salary"`        
	DepartmentID int       `json:"department_id"` 
	Status       string    `json:"status"`        
	CreatedAt    time.Time `json:"created_at"`    
}

// Establish connection to MySQL database
func connectDB() {
	var err error

	//Open DB connection using DSN from environment
	db, err = sql.Open("mysql", os.Getenv("MYSQL_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	//Verify database connectivity
	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable")
	}
	log.Println("Connected to MySQL")
}

// @Summary Create Department
// @Tags Departments
// @Accept json
// @Produce json
// @Param department body Department true "Department Data"
// @Success 200 {object} Department
// @Failure 400 {string} string
// @Router /departments [post]
// Create a new department
func createDepartment(w http.ResponseWriter, r *http.Request) {

	var dept Department
	//Decode request JSON body
	err := json.NewDecoder(r.Body).Decode(&dept)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//Validate department data
	if err := validateDepartment(dept); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//INSERT department into database
	query := "INSERT INTO departments(name) VALUES(?)"
	result, err := db.Exec(query, dept.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get generated ID
	id, _ := result.LastInsertId()
	dept.ID = int(id)
	//Send response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dept)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get Departments
// @Tags Departments
// @Produce json
// @Success 200 {array} Department
// @Router /departments [get]
// Fetch all departments
func getDepartments(w http.ResponseWriter, r *http.Request) {
	//Execute SELECT query
	rows, err := db.Query("SELECT id,name FROM departments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	//Iterate over result set
	for rows.Next() {
		var d Department
		rows.Scan(&d.ID, &d.Name)
		departments = append(departments, d)
	}

	//Return departments as JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(departments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Create Employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body Employee true "Employee Data"
// @Success 200 {object} Employee
// @Router /employees [post]
// Create a new Employee
func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	//Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//validate Employee input
	if err := validateEmployee(emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Insert Employee into database
	query := `
	INSERT INTO employees(name,email,phone,salary,department_id,status)
	VALUES (?,?,?,?,?,?)`

	result, err := db.Exec(
		query,
		emp.Name,
		emp.Email,
		emp.Phone,
		emp.Salary,
		emp.DepartmentID,
		emp.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Assign generated ID
	id, _ := result.LastInsertId()
	emp.ID = int(id)

	//Return created Employee
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get Employees
// @Tags Employees
// @Produce json
// @Success 200 {array} Employee
// @Router /employees [get]
// Fetch all Employees
func getEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
	SELECT id,name,email,phone,salary,department_id,status,created_at
	FROM employees`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	//Map DB rows to struct
	for rows.Next() {
		var e Employee
		rows.Scan(
			&e.ID,
			&e.Name,
			&e.Email,
			&e.Phone,
			&e.Salary,
			&e.DepartmentID,
			&e.Status,
			&e.CreatedAt,
		)
		employees = append(employees, e)
	}
	//Return Employees list
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get Employee by ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} Employee
// @Failure 404 {string} string
// @Router /employees/{id} [get]
// Fetch Employee by ID
func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	//Read ID from URL path
	id := mux.Vars(r)["id"]

	var e Employee
	query := `
	SELECT id,name,email,phone,salary,department_id,status,created_at
	FROM employees WHERE id=?`

	//Query employee
	err := db.QueryRow(query, id).Scan(
		&e.ID,
		&e.Name,
		&e.Email,
		&e.Phone,
		&e.Salary,
		&e.DepartmentID,
		&e.Status,
		&e.CreatedAt,
	)
	//handle record not found
	if err != nil {
		http.Error(w, fmt.Sprintf("Employee not found: %v", err), http.StatusNotFound)
		return
	}
	//return employee
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update Employee
// @Tags Employees
// @Accept json
// @Param id path int true "Employee ID"
// @Param employee body Employee true "Employee Data"
// @Success 200 {string} string
// @Router /employees/{id} [put]
// Update existing employee
func updateEmployee(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var emp Employee
	//decode JSON body
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	//validate Updated data
	if err := validateEmployee(emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//Execute update query
	query := `
	UPDATE employees SET 
	name=?, email=?, phone=?, salary=?, department_id=?, status=?
	WHERE id=?`

	_, err1 := db.Exec(
		query,
		emp.Name,
		emp.Email,
		emp.Phone,
		emp.Salary,
		emp.DepartmentID,
		emp.Status,
		id,
	)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Employee updated successfully."))
}

// @Summary Delete Employee
// @Tags Employees
// @Param id path int true "Employee ID"
// @Success 200 {string} string
// @Router /employees/{id} [delete]
// Delete employee by ID
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	//Execute delete query
	_, err := db.Exec("DELETE FROM employees WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Employee Deleted successfully"))
}

// Validate department fields
func validateDepartment(d Department) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("Department name is required.")
	}
	if len(d.Name) < 2 {
		return errors.New("Department name must be atleast 2 characters")
	}
	return nil
}

// validate employee fields
func validateEmployee(e Employee) error {
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("Employee name is required")
	}
	if !strings.Contains(e.Email, "@") {
		return errors.New("Invalid ID Format")
	}
	if !strings.HasSuffix(e.Email, "@gmail.com") {
		return errors.New("invalid email and does not contain @gmail.com")
	}
	prefix := strings.TrimSuffix(e.Email, "@gmail.com")
	if prefix == "" {
		return errors.New("email must contain prefix to @gmail.com")
	}
	if strings.TrimSpace(e.Phone) == "" {
		return errors.New("phone number is required")
	}
	if e.Salary < 0 {
		return errors.New("salary must be greater than Zero")
	}
	if e.DepartmentID < 0 {
		return errors.New("valid department_id is required")
	}
	if strings.TrimSpace(e.Status) == "" {
		return errors.New("status is required")
	}
	return nil
}

// @title Employee Management API
// @version 1.0
// @description REST API for managing departments and employees
// @host localhost:8080
// @BasePath /
// Application entry point
func Handler() {
	//Load environment variables
	godotenv.Load()

	//connect to DB
	connectDB()
	defer db.Close()

	//setUp router
	router := mux.NewRouter()
	//swagger routes
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	//department routes
	router.HandleFunc("/departments", createDepartment).Methods("POST")
	router.HandleFunc("/departments", getDepartments).Methods("GET")

	//employee routes
	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", getEmployeeByID).Methods("GET")
	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	log.Println("Server running on: 8080")
	http.ListenAndServe(":8080", router)
}
