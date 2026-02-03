package project

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func connectDB() {
	var err error
	dsn := "root:root@tcp(127.0.0.1:3306)/emp_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable")
	}
	log.Println("Connected to MySQL")
}

func createDepartment(w http.ResponseWriter, r *http.Request) {

	var dept Department
	err := json.NewDecoder(r.Body).Decode(&dept)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO departments(name) VALUES(?)"
	result, err := db.Exec(query, dept.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	dept.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(dept)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getDepartments(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id,name FROM departments")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var departments []Department
	for rows.Next() {
		var d Department
		rows.Scan(&d.ID, &d.Name)
		departments = append(departments, d)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(departments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO employees(name,email,phone,salary,department_id,status)
	VALUES (?,?,?,?,?,?)`

	result, err := db.Exec(
		query,
		emp.Phone,
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
	id, _ := result.LastInsertId()
	emp.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(employees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var e Employee
	query := `
	SELECT id,name,email,phone,salary,department_id,status,created_at
	FROM employees WHERE id=?`

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
	if err != nil {
		http.Error(w, "Employee not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var emp Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
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

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM employees WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Employee Deleted successfully"))
}

func Handler() {
	connectDB()
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/departments", createDepartment).Methods("POST")
	router.HandleFunc("/departments", getDepartments).Methods("GET")

	router.HandleFunc("/employees", createEmployee).Methods("POST")
	router.HandleFunc("/employees", getEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", getEmployeeByID).Methods("GET")
	router.HandleFunc("/employees/{id}", updateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", deleteEmployee).Methods("DELETE")

	log.Println("Server running on: 8080")
	http.ListenAndServe(":8080", router)
}
