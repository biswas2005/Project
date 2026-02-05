# Employee Management System

A simple RESTful API built in **Go (Golang)** to manage employee data with **CRUD functionality** and persistence using a SQL database.  
This project provides basic API endpoints to create, read, update, and delete employee records — suitable as a backend for learning, integration, or further expansion.
***

##  Overview

This microservice implements:
- A REST API for employee data
- Database access via Go
- A modular MVC-style code organization
- Migration scripts to create the database schema

---

##  Tech Stack

- **Go** — backend logic  
- **MySQL** (or any SQL-based DB) — persistent storage  
- **Go modules** for dependency management  
- Simple routing via Go’s `net/http` or a router of your choice  

---

##  Project Structure

├── db/migrations/       # SQL migration files to set up schema  
├── docs/                # Project documentation (Swagger)  
├── project/             # Core application code  
├── main.go              # Application entrypoint  
├── go.mod               # Module and dependencies  
├── go.sum               # Dependency versions  
├── .gitignore  
└── README.md
***

##  Prerequisites

Before you begin, make sure you have:

-  Go (version 1.18+)  
-  MySQL server (or other SQL DB)  
-  Git  

---

##  Database Setup

1. **Create the database**
```sql
CREATE DATABASE IF NOT EXISTS emp_db;
````

2. **Run migration SQL files**
   Navigate to `db/migrations/` and run the SQL scripts manually in your SQL client:

```bash
mysql -u root -p emp_db < db/migrations/001_create_employees.sql
```

> Replace `001_create_employees.sql` with the actual migration filenames in your repo.

---

##  Configuration

Create a `.env` file or use environment variables such as:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=emp_db
PORT=8080
```

Make sure this file is ignored in `.gitignore`.

---

##  Run the Application

From your project root:

```bash
go mod tidy
go run main.go
```

By default, your API will run at:

```
http://localhost:8080
```

---

##  API Endpoints

| Method | Endpoint          | Description                 |  
| ------ | ----------------- | --------------------------- |  
| GET    | `/employees`      | Get all employees           |  
| GET    | `/employees/{id}` | Get a specific employee     |  
| POST   | `/employees`      | Create a new employee       |  
| PUT    | `/employees/{id}` | Update an existing employee |  
| DELETE | `/employees/{id}` | Delete an employee          |  

---

##  Example Requests

### Create an employee

```bash
curl -X POST http://localhost:8080/employees ^
-H "Content-Type: application/json" 
-d "{\"name\":\"Alice Johnson\",\"email\":\"alice@example.com\",\"phone\":\"9876543210\",\"salary\":75000,\"department_id\":1,\"status\":\"Active\"}"

```

### Get all employees

```bash
curl http://localhost:8080/employees
```
### Get employee by ID
```bash
curl http://localhost:8080/employee/1
```

### Update an employee

```bash
curl -X PUT http://localhost:8080/employees/1 ^
-H "Content-Type: application/json" ^
-d "{\"name\":\"Mark Henry\",\"email\":\"henry@example.com\",\"phone\":\"0123456780\",\"salary\":75000,\"department_id\":1,\"status\":\"Active\"}"

```

### Delete an employee

```bash
curl -X DELETE http://localhost:8080/employees/1
```

---

##  Status Codes

| Code | Meaning      |
| ---- | ------------ |
| 200  | OK           |
| 201  | Created      |
| 400  | Bad Request  |
| 404  | Not Found    |
| 500  | Server Error |

---


##  Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch
3. Commit changes
4. Push and open a Pull Request

---

##  License

This project is licensed under **Apache-2.0 License**. ([GitHub][1])

---

##  Thanks

Thank you for visiting the Employee Management System repository.
Feel free to reach out if you want help setting up, extending, or deploying this project!




